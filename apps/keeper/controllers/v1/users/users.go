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
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
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

// RenewAccessToken
// @Summary Renew access token
// @Tags    Users
// @Accept  json
// @Produce json
// @Success 200 {object} utils.ResponseOk
// @Failure 400 {object} utils.ResponseError
// @Router  /v1/users/renew-token
func (u UsersController) RenewAccessToken(c *gin.Context) {
	refreshToken, err := c.Request.Cookie("Refresh-Token")
	if err != nil {
		errorsEncountered := utils.GenerateValidationErrors(errors.New("no refresh token was present in cookies"))
		c.JSON(400, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Find the existing refresh token
	rto, err := u.RefreshTokenStorage.FindExistingRefreshToken(refreshToken.Value)
	if err != nil {
		errorsEncountered := utils.GenerateValidationErrors(errors.New("refresh token is invalid"))
		c.JSON(404, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	currentRefreshTokenExpiresIn := rto.Expires
	currentTime := time.Now()

	// If this refresh token is already expired then log the user out
	if currentTime.After(*currentRefreshTokenExpiresIn) {
		tempExpirationTime := time.Now()
		errorsEncountered := utils.GenerateValidationErrors(errors.New("refresh token is no longer valid, please login again"))
		c.SetCookie("Access-Token", "", int(tempExpirationTime.UnixMilli()), "/", os.Getenv("JWT_COOKIE_DOMAIN"), false, true)
		c.JSON(401, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	user, err := u.Storage.FindUserWithUID(rto.UserID)
	if err != nil {
		errorsEncountered := utils.GenerateErrorResponse(err)
		c.JSON(500, utils.ResponseError{Errors: errorsEncountered, IsError: true})
		return
	}

	// Generate a new access token for the user
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Email:    user.Email,
		UserID:   user.UserID,
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

	c.SetCookie("Access-Token", accessTokenString, int(expirationTime.UnixMilli()), "/", os.Getenv("JWT_COOKIE_DOMAIN"), false, true)
	c.JSON(200, utils.ResponseOk{Message: "Logged in successfully!"})
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
		Email:    existingUser.Email,
		UserID:   existingUser.UserID,
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

// Logout
// @Summary  Logs the user out
// @Tags     Users
// @Accept   json
// @Produce  json
// @Success  200  {object}  utils.ResponseOk
// @Failure  500  {object}  utils.ResponseError
// @Router   /v1/users/logout [get]
func (u UsersController) LogoutUser(c *gin.Context) {
	// Generate a new access token for the user
	expirationTime := time.Now()
	claims := &Claims{
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

	c.SetCookie("Access-Token", accessTokenString, int(expirationTime.UnixMilli()), "/", os.Getenv("JWT_COOKIE_DOMAIN"), false, true)
	c.JSON(200, utils.ResponseOk{Message: "Logged out successfully!"})
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
