package routes

import (
	users_v1 "apps/keeper/controllers/v1/users"
	"apps/keeper/utils"
)

// UsersRoutes struct
type UsersRoutes struct {
	logger          utils.Logger
	handler         utils.RequestHandler
	usersController users_v1.UsersController
}

// Setup Users routes
func (s UsersRoutes) Setup() {
	s.logger.Info("Setting up routes")
	api := s.handler.Gin.Group("/apis/v1/users")
	{
		api.POST("/register", s.usersController.CreateUser)
		api.POST("/login", s.usersController.LoginUser)
	}
}

// NewUsersRoutes creates new Users controller
func NewUsersRoutes(
	logger utils.Logger,
	handler utils.RequestHandler,
	usersController users_v1.UsersController,
) UsersRoutes {
	return UsersRoutes{
		handler:         handler,
		logger:          logger,
		usersController: usersController,
	}
}
