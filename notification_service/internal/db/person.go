package db

import (
	"context"
	"fmt"
	"notification_service/internal/models"
)

func CreatePerson(person models.PersonCreate, conn *PostgresConnection, ctx context.Context) (int, error) {
	queryString := `
INSERT INTO persons (email, telegram_id, phone_number)
VALUES ($1, $2, $3)
RETURNING person_id;
	`
	row := conn.DB.QueryRowContext(
		ctx,
		queryString,
		person.Email.Raw,
		person.TelegramId.Raw,
		person.PhoneNumber.Raw,
	)
	var personId int
	err := row.Scan(&personId)
	if err != nil {
		return 0, err
	}

	return personId, nil
}

func DeletePersons(conn *PostgresConnection, ctx context.Context, personIds []int) error {
	// Building query string
	queryString := "DELETE FROM persons WHERE person_id IN ("

	args := make([]any, len(personIds))
	for idx, id := range personIds {
		queryString += fmt.Sprintf("$%d", idx+1)
		if idx != len(personIds)-1 {
			queryString += ", "
		}
		args[idx] = id
	}
	queryString += ");"

	// Doing query
	_, err := conn.DB.ExecContext(
		ctx,
		queryString,
		args...,
	)

	return err
}

func GetAllPersons(conn *PostgresConnection, ctx context.Context) ([]models.Person, error) {
	queryString := `
SELECT * FROM persons;
	`
	rows, err := conn.DB.QueryContext(ctx, queryString)
	if err != nil {
		return []models.Person{}, nil
	}

	result := make([]models.Person, 0)
	for rows.Next() {
		curValue := models.Person{}
		rows.Scan(&curValue.PersonId, &curValue.Email.Raw, &curValue.TelegramId.Raw, &curValue.PhoneNumber.Raw)
		result = append(result, curValue)
	}

	return result, nil
}

func GetPerson(conn *PostgresConnection, ctx context.Context, personId int) (models.Person, error) {
	queryString := `
SELECT * FROM persons WHERE person_id = $1;
	`
	row := conn.DB.QueryRowContext(ctx, queryString, personId)
	var person models.Person
	if err := row.Scan(
		&person.PersonId,
		&person.Email.Raw,
		&person.TelegramId.Raw,
		&person.PhoneNumber.Raw,
	); err != nil {
		return models.Person{}, err
	}
	return person, nil
}

func GetPersonsByGroupsFilter(conn *PostgresConnection, ctx context.Context, groupNames []string) ([]models.Person, error) {
	if len(groupNames) == 0 {
		return make([]models.Person, 0), nil
	}

	placeholders := "("
	args := make([]any, len(groupNames))
	for idx, name := range groupNames {
		args[idx] = name
		placeholders += fmt.Sprintf("$%d", idx+1)
		if idx != len(groupNames)-1 {
			placeholders += ", "
		}
	}
	placeholders += ")"

	queryString := fmt.Sprintf(`
WITH selected_group_ids AS (
	SELECT group_id
	FROM groups
	WHERE name IN %s 
)
SELECT DISTINCT p.person_id, p.email, p.telegram_id, p.phone_number
FROM person_to_group AS ptg 
	JOIN persons AS p
	ON ptg.person_id = p.person_id
WHERE   
	ptg.group_id IN (SELECT * FROM selected_group_ids);
	`, placeholders)

	rows, err := conn.DB.QueryContext(ctx, queryString, args...)
	if err != nil {
		return []models.Person{}, err
	}

	result := make([]models.Person, 0)
	for rows.Next() {
		var temp models.Person
		if err := rows.Scan(&temp.PersonId, &temp.Email.Raw, &temp.TelegramId.Raw, &temp.PhoneNumber.Raw); err != nil {
			return []models.Person{}, err
		}
		result = append(result, temp)
	}
	return result, nil
}
