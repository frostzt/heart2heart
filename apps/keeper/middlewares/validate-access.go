package middlewares

import (
	"apps/keeper/storage"
	"apps/keeper/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

// This will be encoded in the JWT exchanged
type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

// Validate Access Token sent by the user and validate if the user has
// access to the resource or not
func ValidateAccessToken(fn gin.HandlerFunc, userStorage storage.UserStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from the request
		accessToken, err := c.Request.Cookie("Access-Token")
		if err != nil {
			c.JSON(401, utils.ErrorMessage{Message: "no access token provided please login first"})
			return
		}

		// this also checks internally for claims including token expiry
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(accessToken.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return JWT_SECRET, nil
		})

		// Something failed while getting the JWT_SECRET, or something went wrong with parsing
		if err != nil {
			c.JSON(401, utils.ErrorMessage{Message: "token invalid or expired please relogin"})
			return
		}

		// Check if the provided token is valid or not
		if !tkn.Valid {
			c.JSON(401, utils.ErrorMessage{Message: "your session is invalid please relogin"})
			return
		}

		// Set the validated claims in the gin context for the handler
		c.Set("Validated-Claims", claims)
		fn(c)
	}
}
