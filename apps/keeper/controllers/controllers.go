package controllers

import (
	misc_v1 "apps/keeper/controllers/v1/misc"
	users_v1 "apps/keeper/controllers/v1/users"

	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(misc_v1.NewMiscController),
	fx.Provide(users_v1.NewUsersController),
)
