package db

import (
	"database/sql"
	"fmt"
	"notification_service/internal/config"

	_ "github.com/lib/pq"
)

type PostgresConnection struct {
	DB *sql.DB
}

func CreatePostgresConnection(config *config.DatabaseConfig) (PostgresConnection, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DbName,
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
