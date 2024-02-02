package connect

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var dbUrl = "postgresql://postgres:postgres@localhost:5432/aqary-db"

var initialize = false

func ConnectDB(*ctx context.Context) *pgx.Conn {
	conn, err := pgx.Connect(ctx)
	if err != nil {
		panic(err)
	}

	if initialize {
		schema, err := os.ReadFile("models/schema/schema.sql")
		if err != nil {
			log.Fatalf("Error reading schema.sql: %v", err)
		}
		_, err = conn.Exec(ctx, string(schema))
		if err != nil {
			log.Fatalf("Error executing schema.sql: %v", err)
		}
		defer conn.Close(ctx)
	}

	return conn
}
