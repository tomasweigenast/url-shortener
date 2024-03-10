package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

// var connection *pgx.Conn
var pool *pgxpool.Pool

func ConnectDatabase() {
	dbUrl, ok := os.LookupEnv("DATABASE_URL")

	if !ok {
		log.Fatalf("DATABASE_URL environment variable is not found")
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbUrl)
	// connection, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("could not connect to database: %s", err)
	}

	// try connection
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("could not reach database: %s", err)
	}

	prepareDb()
}

func CloseDatabase() {
	if pool != nil {
		pool.Close()
		pool = nil
	}
}

func Pool() *pgxpool.Pool {
	return pool
}

func prepareDb() {
	_, err := pool.Exec(context.Background(), queryCreateTables)
	if err != nil {
		log.Fatalf("unable to execute create table script, error: %s", err)
	}

	migrator := rivermigrate.New(riverpgxv5.New(pool), nil)
	_, err = migrator.Migrate(context.Background(), rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{})
	if err != nil {
		log.Fatalf("unable to migrate db for river: %s", err)
	}
}
