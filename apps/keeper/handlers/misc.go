package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Healthcheck godoc
// @Summary Healthcheck API for the keeper service
// @Description Returns 200 if the API is up and running
// @ID health-check
// @Tags misc
// @Produce  text
// @Success 200 {object} ResponseOk
// @Failure 500 {object}
// @Router /healthcheck [get]
func (h *Handler) Healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, ResponseOk{Message: "Keeper API accepting connections..."})
}
