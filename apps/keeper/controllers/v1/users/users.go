package users_v1

import (
	"apps/keeper/storage"
	"apps/keeper/utils"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const VERSION string = "development"

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

// This will be encoded in the JWT exchanged
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Struct validator
var validate *validator.Validate

// UsersController
type UsersController struct {
	logger              utils.Logger
	Storage             *storage.UserStorage
	RefreshTokenStorage *storage.RefreshTokenStorage
}

// NewUsersController creates new Users Controller
func NewUsersController(storage *storage.UserStorage, refreshTokenStorage *storage.RefreshTokenStorage, logger utils.Logger) UsersController {
	return UsersController{
		logger:              logger,
		Storage:             storage,
		RefreshTokenStorage: refreshTokenStorage,
	}
}

// LoginUser
// @Summary Log the user in
// @Tags    Users
// @Accept  json
// @Produce json
// @Success 200 {object} utils.ResponseOk
// @Failure 400 {object} utils.ResponseError
// @Router  /v1/users/login
func (u UsersController) LoginUser(c *gin.Context) {
	var user storage.LoginUser

	// Validations
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := c.BindJSON(&user)

	validationError := validate.Struct(user)
	if validationError != nil {
		errorsEncountered := utils.GenerateValidationErrors(validationError)
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(err)
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Logic
	existingUser, err := u.Storage.FindUserWithUsername(user.Username, true)
	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(err)
		c.JSON(500, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// User not found
	if existingUser == nil {
		errorsEncountered := utils.GenerateErrorResponse(errors.New("invalid credentials"))
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Invalid password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(errors.New("invalid credentials"))
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Valid password was found, generate access token
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString(JWT_SECRET)
	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(err)
		c.JSON(500, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Generate long-lived refresh token
	refreshTokenObject, err := u.RefreshTokenStorage.GenerateAndInsertRefreshToken(existingUser)
	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(err)
		c.JSON(500, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	c.SetCookie("Access-Token", accessTokenString, int(expirationTime.UnixMilli()), "/", os.Getenv("JWT_COOKIE_DOMAIN"), false, true)
	c.SetCookie("Refresh-Token", refreshTokenObject.Token, int(refreshTokenObject.Expires.UnixMilli()), "/", os.Getenv("JWT_COOKIE_DOMAIN"), false, true)

	c.JSON(200, utils.ResponseOk{Message: "Logged in successfully!"})
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
	var newUser storage.UserData

	// Validations
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := c.BindJSON(&newUser)

	validationError := validate.Struct(newUser)
	if validationError != nil {
		errorsEncountered := utils.GenerateValidationErrors(validationError)
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Check for binding error if one was received
	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(err)
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered})
		return
	}

	// Logic
	responseCode, err := u.Storage.CreateNewUser(storage.UserData{
		Name:      newUser.Name,
		BirthDate: newUser.BirthDate,
		Username:  newUser.Username,
		Email:     newUser.Email,
		Password:  newUser.Password,
	})

	if err != nil && responseCode != 0 {
		c.JSON(500, utils.ErrorMessage{Message: err.Error()})
		return
	}

	c.JSON(202, utils.ResponseOk{Message: "User made successfully, please login!"})
}
