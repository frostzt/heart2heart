package users_v1

import (
	"apps/keeper/storage"
	"apps/keeper/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const VERSION string = "development"

// Struct validator
var validate *validator.Validate

// UsersController
type UsersController struct {
	logger  utils.Logger
	Storage *storage.UserStorage
}

// Data that will be received from the user
type UserData struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// NewUsersController creates new Users Controller
func NewUsersController(storage *storage.UserStorage, logger utils.Logger) UsersController {
	return UsersController{
		logger:  logger,
		Storage: storage,
	}
}

// CreateUser
// @Summary  Create a new user
// @Tags     Users
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/users/register [post]
func (u UsersController) CreateUser(c *gin.Context) {
	var newUser UserData

	validate = validator.New(validator.WithRequiredStructEnabled())
	err := c.BindJSON(&newUser)

	// Validate struct that was recieved from the user
	validationError := validate.Struct(newUser)
	if validationError != nil {
		errorsEncountered := []utils.ErrorMessage{}

		for _, e := range validationError.(validator.ValidationErrors) {
			errorMessage := utils.ErrorMessage{
				Message: e.Error(),
			}

			errorsEncountered = append(errorsEncountered, errorMessage)
		}

		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Check for binding error if one was received
	if err != nil {
		errorMessage := utils.ErrorMessage{Message: err.Error()}
		errorsEncountered := []utils.ErrorMessage{}
		errorsEncountered = append(errorsEncountered, errorMessage)

		c.JSON(400, utils.ResponseError{Errors: errorsEncountered})
		return
	}

	fizz, err := u.Storage.CreateNewUser(storage.UserData{
		ID:   1,
		Name: "Sourav",
	})

	if err != nil {
		c.JSON(500, utils.ErrorMessage{Message: err.Error()})
		return
	}

	fmt.Println(fizz)

	c.JSON(202, utils.ResponseOk{Message: "User made!"})
}
