package handlers

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) Register(v1 *echo.Group) {
	// Users specific routes
	users := v1.Group("/users")

	users.POST("/register", h.SignUp)
	users.POST("/login", h.LoginUser)
}
