package controllers

import (
	misc_v1 "apps/keeper/controllers/v1/misc"
	"go.uber.org/fx"
)

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(misc_v1.NewMiscController),
)
