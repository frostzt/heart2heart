package database

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type DBParams struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func CreatePostgresConnection() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		fmt.Println("💥 Wasn't able to parse provided port for db, defaulting to 5432!")
		dbPort = 5432
	}

	connString := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s", dbHost, dbPort, dbName, dbUser, dbPass)

	// Create the connection pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	fmt.Println("🛅 Connected to Postgres...")

	return pool, nil
}
