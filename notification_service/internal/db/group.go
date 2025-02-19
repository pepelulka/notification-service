package db

import (
	"context"
	"fmt"
	"log"
	"notification_service/internal/models"
)

func GetAllGroupNames(conn *PostgresConnection, ctx context.Context) ([]string, error) {
	queryString := `
SELECT name FROM groups;
	`

	rows, err := conn.DB.QueryContext(ctx, queryString)
	if err != nil {
		return []string{}, err
	}

	result := make([]string, 0)
	for rows.Next() {
		var cur string
		if err = rows.Scan(&cur); err != nil {
			return []string{}, err
		}
		result = append(result, cur)
	}
	return result, nil
}

func CreateGroup(group models.GroupCreate, conn *PostgresConnection, ctx context.Context) error {
	tx, err := conn.DB.BeginTx(ctx, nil) // Start transaction
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(
		ctx,
		"INSERT INTO groups (name) VALUES ($1) RETURNING group_id;",
		group.Name,
	)

	var groupId int
	if err := row.Scan(&groupId); err != nil {
		if rollbackError := tx.Rollback(); rollbackError != nil {
			log.Fatalf("rollback error: %v", rollbackError)
		}
		return err
	}

	// Start to build big query
	if len(group.ParticipantIds) != 0 {
		placeholders := ""
		query := "INSERT INTO person_to_group (person_id, group_id) VALUES "
		args := make([]any, len(group.ParticipantIds))
		for idx, person_id := range group.ParticipantIds {
			placeholders += fmt.Sprintf("($%d, %d)", idx+1, groupId)
			if idx != len(group.ParticipantIds)-1 {
				placeholders += ", "
			}
			args[idx] = person_id
		}
		query += placeholders

		_, err := tx.ExecContext(
			ctx,
			query,
			args...,
		)

		if err != nil {
			if rollbackError := tx.Rollback(); rollbackError != nil {
				log.Fatalf("rollback error: %v", rollbackError)
			}
			return err
		}
	}

	return tx.Commit()
}

func AddPersonsToGroup(groupName string, personIds []int, conn *PostgresConnection, ctx context.Context) error {
	if len(personIds) == 0 {
		return nil
	}
	queryString := `
WITH given_group_id AS (
	SELECT group_id 
	FROM groups
	WHERE name = $1
)
INSERT INTO person_to_group (person_id, group_id)
VALUES 
	`
	args := make([]any, len(personIds)+1)
	args[0] = groupName
	for idx, personId := range personIds {
		queryString += fmt.Sprintf("($%d, (SELECT * FROM given_group_id) )", idx+2)
		if idx != len(personIds)-1 {
			queryString += ", "
		}
		args[idx+1] = personId
	}
	queryString += " ON CONFLICT DO NOTHING;"
	_, err := conn.DB.ExecContext(ctx, queryString, args...)
	return err
}

func GetGroupByName(groupName string, conn *PostgresConnection, ctx context.Context) (models.Group, error) {
	// At first we need to get group id
	row := conn.DB.QueryRowContext(
		ctx,
		"SELECT group_id FROM groups WHERE name = $1;",
		groupName,
	)

	var groupId int
	if err := row.Scan(&groupId); err != nil {
		return models.Group{}, err
	}

	// At second we need to get all relations to this group
	rows, err := conn.DB.QueryContext(
		ctx,
		"SELECT person_id FROM person_to_group WHERE group_id = $1;",
		groupId,
	)
	if err != nil {
		return models.Group{}, err
	}

	participantIds := make([]int, 0)
	for rows.Next() {
		var temp int
		if err := rows.Scan(&temp); err != nil {
			return models.Group{}, err
		}
		participantIds = append(participantIds, temp)
	}

	return models.Group{
		GroupId:        groupId,
		Name:           groupName,
		ParticipantIds: participantIds,
	}, nil
}

func DeleteGroupByName(name string, conn *PostgresConnection, ctx context.Context) error {
	_, err := conn.DB.ExecContext(
		ctx,
		"DELETE FROM groups WHERE name = $1;",
		name,
	)
	return err
}
