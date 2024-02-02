package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitDB() (*pgx.Conn, error) {
	databaseURL := os.Getenv("DB_URL")
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}
	if err := Init(conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func Init(conn *pgx.Conn) error {
	initSQL, err := os.ReadFile("sql/schema.sql")
	if err != nil {
		log.Printf("Error while initializing the table: %v", err)
		return err
	}

	_, err = conn.Exec(context.Background(), string(initSQL))

	if err != nil {
		log.Printf("Error on execution of table creation: %v", err)
		return err
	}
	log.Println("Database table is created successfully")
	return nil
}
