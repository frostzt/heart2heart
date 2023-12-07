package handlers

import (
	"errors"
	"fmt"
	"keeper/storage"
	"keeper/utils"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// This will be encoded in the JWT exchanged
type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

// Struct validator
var validate *validator.Validate

// LoginUser godoc
// @Summary Log an existing user in
// @Description Log an existing user in
// @ID login-user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body storage.LoginUser true "User info for registration"
// @Success 200 {object} ResponseOk
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /login [post]
func (h *Handler) LoginUser(c echo.Context) error {
	var user storage.LoginUser

	// Validations
	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	validationError := validate.Struct(user)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, utils.NewValidatorError(validationError))
	}

	// Logic
	existingUser, err := h.userStorage.FindUserWithUsername(user.Username, true)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("invalid credentials")))
		}

		fmt.Println("💥 Error encountered while fetching user:", err)
		return c.JSON(http.StatusInternalServerError, utils.NewError(err))
	}

	// Invalid password
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(errors.New("invalid credentials")))
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
		fmt.Println("💥 Error encountered while signing token:", err)
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("something went wrong")))
	}

	// Generate long-lived refresh token
	refreshTokenObject, err := h.refreshTokenStorage.GenerateAndInsertRefreshToken(existingUser)
	if err != nil {
		fmt.Println("💥 Error encountered while generating refresh token:", err)
		return c.JSON(http.StatusInternalServerError, utils.NewError(errors.New("something went wrong")))
	}

	accessTokenCookie := utils.GenerateCookie("Access-Token", accessTokenString, os.Getenv("JWT_COOKIE_DOMAIN"), "/", true, expirationTime)
	refreshTokenCookie := utils.GenerateCookie("Refresh-Token", refreshTokenObject.Token, os.Getenv("JWT_COOKIE_DOMAIN"), "/", true, *refreshTokenObject.Expires)

	c.SetCookie(&accessTokenCookie)
	c.SetCookie(&refreshTokenCookie)

	return c.JSON(http.StatusOK, ResponseOk{Message: "Logged in successfully!"})
}

// SignUp godoc
// @Summary Register a new user
// @Description Register a new user
// @ID sign-up
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body storage.UserData true "User info for registration"
// @Success 201 {object} ResponseOk
// @Failure 400 {object} utils.Error
// @Failure 404 {object} utils.Error
// @Failure 500 {object} utils.Error
// @Router /users [post]
func (h *Handler) SignUp(c echo.Context) error {
	var newUser storage.UserData

	// Validations
	validate = validator.New(validator.WithRequiredStructEnabled())
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	validationError := validate.Struct(newUser)
	if validationError != nil {
		return c.JSON(http.StatusBadRequest, utils.NewValidatorError(validationError))
	}

	// Create the user
	_, err := h.userStorage.CreateNewUser(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewError(err))
	}

	return c.JSON(http.StatusCreated, ResponseOk{Message: "registered successfully please login"})
}
