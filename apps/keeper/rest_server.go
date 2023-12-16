package main

import (
	"fmt"
	"keeper/db"
	"keeper/handlers"
	"keeper/router"
	"keeper/storage"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func StartRESTServer() {
	r := router.New()

	// Setup swagger - we're using Swaggo
	r.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := r.Group("/v1/api")

	// Establish Connection with Database
	conn, err := db.CreatePostgresConnection()
	if err != nil {
		fmt.Println("💣 Error encountered while connecting to database: ", err.Error())
		return
	}

	// Migrate the database
	err = MigrateDatabaseOnStart()
	if err != nil {
		fmt.Println("💣 Error encountered while migrating database:", err.Error())
		return
	}

	us := storage.NewUserStore(conn)
	rts := storage.NewRefreshTokenStorage(conn)

	h := handlers.NewHandler(us, rts)
	h.Register(v1)

	// Init REST API Server
	r.Logger.Fatal(r.Start("0.0.0.0:1323"))
}

func MigrateDatabaseOnStart() error {
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		fmt.Println("💥 Wasn't able to parse provided port for db, defaulting to 5432!")
		dbPort = 5432
	}

	// Run migrations
	dbUri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	m, err := migrate.New("file:///app/db/migrations", dbUri)

	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	fmt.Println("✅ Database successfully migrated...")
	return nil
}
