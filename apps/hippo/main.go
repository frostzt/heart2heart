package main

import (
    "github.com/joho/godotenv"
	"go.uber.org/fx"
	"apps/hippo/utils"
	"apps/hippo/bootstrap"
)

func main() {
	godotenv.Load()
	logger := utils.GetLogger().GetFxLogger()
	fx.New(bootstrap.Module, fx.Logger(logger)).Run()
}
