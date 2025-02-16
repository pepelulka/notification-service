package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"notification_service/internal/config"
	"notification_service/internal/models"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	DB *sql.DB
}

func CreatePostgresConnection(config config.Config) (PostgresConnection, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DbName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return PostgresConnection{}, err
	}
	return PostgresConnection{
		DB: db,
	}, nil
}

func (conn *PostgresConnection) Close() error {
	return conn.DB.Close()
}

func CreatePerson(person models.PersonCreate, conn *PostgresConnection, ctx context.Context) error {
	queryString := `
INSERT INTO persons (email, telegram_id, phone_number)
VALUES ($1, $2, $3);
	`
	_, err := conn.DB.ExecContext(
		ctx,
		queryString,
		person.Email,
		person.TelegramId,
		person.PhoneNumber,
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
		args := make([]any, 0)
		for idx, person_id := range group.ParticipantIds {
			placeholders += fmt.Sprintf("($%d, %d)", idx+1, groupId)
			if idx != len(group.ParticipantIds)-1 {
				placeholders += ", "
			}
			args = append(args, person_id)
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

// func (repo *UserRepository) GetInfoByLogin(login string) (models.UserInfo, error) {
// 	var info models.UserInfo

// 	row := repo.DB.QueryRow(
// 		"SELECT login, role, COALESCE(patient_id, 0), COALESCE(doctor_id, 0) "+
// 			"FROM clinic.user_credentials "+
// 			"WHERE login = $1;",
// 		login,
// 	)
// 	if err := row.Scan(&info.Login, &info.Role, &info.PatientId, &info.DoctorId); err != nil {
// 		return models.UserInfo{}, err
// 	}
// 	return info, nil
// }

// func (repo *UserRepository) GetPasswordHashByLogin(login string) (string, error) {
// 	var passwordHash string

// 	row := repo.DB.QueryRow(
// 		"SELECT password_hash "+
// 			"FROM clinic.user_credentials "+
// 			"WHERE login = $1;",
// 		login,
// 	)
// 	if err := row.Scan(&passwordHash); err != nil {
// 		return "", err
// 	}
// 	return passwordHash, nil
// }

// func (repo *UserRepository) CreateAdmin(user models.AdminCreateHashed) error {
// 	_, err := repo.DB.Exec(
// 		`
// INSERT INTO clinic.user_credentials (login, password_hash, role, patient_id, doctor_id)
// VALUES ($1, $2, 'admin', NULL, NULL);
// 		`,
// 		user.Login,
// 		user.PasswordHash,
// 	)
// 	return err
// }

// func (repo *UserRepository) CreatePatient(patient models.PatientCreateHashed) error {
// 	_, err := repo.DB.Exec(
// 		`
// INSERT INTO clinic.user_credentials (login, password_hash, role, patient_id, doctor_id)
// VALUES ($1, $2, 'patient', $3, NULL);
// 		`,
// 		patient.Login,
// 		patient.PasswordHash,
// 		patient.PatientId,
// 	)
// 	return err
// }

// func (repo *UserRepository) CreateDoctor(patient models.DoctorCreateHashed) error {
// 	_, err := repo.DB.Exec(
// 		`
// INSERT INTO clinic.user_credentials (login, password_hash, role, patient_id, doctor_id)
// VALUES ($1, $2, 'doctor', NULL, $3);
// 		`,
// 		patient.Login,
// 		patient.PasswordHash,
// 		patient.DoctorId,
// 	)
// 	return err
// }
