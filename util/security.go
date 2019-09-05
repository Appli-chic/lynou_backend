package util

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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
		return []byte("AllYourBase"), nil
	})

	// Check if the token is correct and valid
	if err != nil || token == nil || !token.Valid {
		return nil, errors.New("no token found")
	}

	return token, nil
}

// Retrieve the token to check if the user is authenticated
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
