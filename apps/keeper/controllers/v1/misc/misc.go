package misc_v1

import (
	"apps/keeper/utils"

	"github.com/gin-gonic/gin"
)

// use ldflags to replace this value during build:
// 		https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
const VERSION string = "development"

// MiscController data type
type MiscController struct {
	logger          utils.Logger
}

// NewMiscController creates new Misc controller
func NewMiscController(logger utils.Logger) MiscController {
	return MiscController{
		logger:          logger,
	}
}

// GetVersion
// @Summary  Get version of keeper application
// @Tags     Misc
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/version [get]
func (u MiscController) GetVersion(c *gin.Context) {
	c.JSON(200, utils.ResponseOk{Message: VERSION})
}

// GetReadiness
// @Summary  Readiness endpoint
// @Tags     Misc
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/readiness [get]
func (u MiscController) GetReadiness(c *gin.Context) {
  c.JSON(200, utils.ResponseOk{Message: "Keeper API ready..."})
}

// GetLiveness
// @Summary  Liveness endpoint
// @Tags     Misc
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Router   /v1/liveness [get]
func (u MiscController) GetLiveness(c *gin.Context) {
  c.JSON(200, utils.ResponseOk{Message: "Keeper API live..."})
}
