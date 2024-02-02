package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func Init() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))

	if err != nil {
		fmt.Fprint(os.Stderr, "unable to connect to the database")
		os.Exit(1)
	}
	initSql, err := os.ReadFile("sq/schema.sql")
	if err != nil {
		log.Printf("Error while loading the file: %v", err)
		return nil, err
	}

	_, err = conn.Exec(context.Background(), string(initSql))

	if err != nil {
		log.Printf("Error while executing the schema: %v", err)
		return nil, err
	}

	return conn, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Can not load env file!")
	}
	// conn, err := Init()
	// store := database.New(conn)

	router := gin.Default()

	// defer conn.Close(context.Background())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "connected",
		})
	})

}
