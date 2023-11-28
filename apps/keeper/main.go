package main

import (
	"apps/keeper/bootstrap"
	"apps/keeper/utils"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	// Run migrations
	dbUri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	m, err := migrate.New("file:///app/db/migrations", dbUri)

	if err != nil {
		fmt.Println("💣 Error encountered while creating a migrator instance: ", err.Error())
		return
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("💣 Error encountered while migrating to latest: ", err.Error())
		return
	}

	fmt.Println("✅ Database successfully migrated...")

	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
