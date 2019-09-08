package util

import (
	"errors"
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
	"github.com/applichic/lynou/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetUserFromToken(c *gin.Context) (*model.User, error) {
	token, err := GetToken(c)

	if err != nil {
		return nil, err
	}

	userService := new(service.UserService)
	userClaims := token.Claims.(jwt.MapClaims)["User"].(map[string]interface{})
	user, err := userService.FetchUserById(userClaims["ID"])

	return &user, err
}

// Get token from the Authorization header
func GetToken(c *gin.Context) (*jwt.Token, error) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	tokenString := strings.TrimSpace(splitToken[1])

	// Check if there is a token given
	if tokenString == "" {
		return nil, errors.New("no token found")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Conf.JwtSecret), nil
	})

	// Check if the token is correct and valid
	if err != nil || token == nil || !token.Valid {
		return nil, errors.New("no token found")
	}

	return token, nil
}

// Retrieve the token to check if the service is authenticated
func AuthenticationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := GetToken(c)

		// Check if the token is valid
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		return
	}
}
