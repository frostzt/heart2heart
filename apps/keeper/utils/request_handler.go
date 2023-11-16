package utils

import (
	"apps/keeper/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RequestHandler function
type RequestHandler struct {
	Gin *gin.Engine
}

type ErrorMessage struct {
	Message string `json:"message"`
} //@name ErrorMessage

type ResponseError struct {
	IsError bool           `json:"isError"`
	Errors  []ErrorMessage `json:"errors"`
} //@name ResponseError

type ResponseOk struct {
	Message string `json:"message"`
} //@name ResponseOk

func GenerateErrorResponse(err error) []ErrorMessage {
	errorMessage := ErrorMessage{Message: err.Error()}
	errorsEncountered := []ErrorMessage{}
	errorsEncountered = append(errorsEncountered, errorMessage)

	return errorsEncountered
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(logger Logger, env Env) RequestHandler {
	gin.DefaultWriter = logger.GetGinLogger()
	engine := gin.New()
	docs.SwaggerInfo.Title = "keeper"
	docs.SwaggerInfo.BasePath = "/apis"
	engine.GET(
		"/openapi/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
		),
	)

	return RequestHandler{Gin: engine}
}
