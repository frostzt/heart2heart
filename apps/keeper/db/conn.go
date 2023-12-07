package db

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func CreatePostgresConnection() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	// if utils.GetStringAbsoluteLength(dbHost) == 0 {
	// 	return nil, errors.New("'DB_HOST' is not defined")
	// }

	dbName := os.Getenv("DB_NAME")
	// if utils.GetStringAbsoluteLength(dbName) == 0 {
	// 	return nil, errors.New("'DB_NAME' is not defined")
	// }

	dbUser := os.Getenv("DB_USERNAME")
	// if utils.GetStringAbsoluteLength(dbUser) == 0 {
	// 	return nil, errors.New("'DB_USERNAME' is not defined")
	// }

	dbPass := os.Getenv("DB_PASSWORD")
	// if utils.GetStringAbsoluteLength(dbPass) == 0 {
	// 	return nil, errors.New("'DB_PASSWORD' is not defined")
	// }

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

	// _, err = pool.Query(context.Background(), "SELECT 1 + 1;")
	// if err != nil {
	// 	return nil, err
	// }

	fmt.Println("🛅 Connected to Postgres...")

	return pool, nil
}
