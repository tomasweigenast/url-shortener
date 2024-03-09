package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var connection *pgx.Conn

func ConnectDatabase() {
	dbUrl, ok := os.LookupEnv("DATABASE_URL")

	if !ok {
		log.Fatalf("DATABASE_URL environment variable is not found")
	}

	var err error
	connection, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	// try connection
	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatalf("could not reach database: %s", err)
	}

	prepareDb()
}

func CloseDatabase() {
	if connection != nil {
		connection.Close(context.Background())
		connection = nil
	}
}

func GetConnection() *pgx.Conn {
	return connection
}

func prepareDb() {
	_, err := connection.Exec(context.Background(), queryCreateTables)
	if err != nil {
		log.Fatalf("unable to execute create table script, error: %s", err)
	}
}
