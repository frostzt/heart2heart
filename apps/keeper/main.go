package main

import (
    "github.com/joho/godotenv"
	"go.uber.org/fx"
	"apps/keeper/utils"
	"apps/keeper/bootstrap"
)

func main() {
	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
