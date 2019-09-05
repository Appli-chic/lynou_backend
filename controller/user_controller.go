package controller

import (
	"github.com/applichic/lynou/database"
	"github.com/applichic/lynou/model"
	validator2 "github.com/applichic/lynou/validator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"time"
)

const codeErrorServer = "CODE_ERROR_SERVER"
const codeErrorEmailAlreadyExists = "CODE_ERROR_EMAIL_ALREADY_EXISTS"

var jwtKey = []byte("my_secret_key")
var bcryptSecret = []byte("MyDarkSecret")

type UserController struct {
}

// Sign up the user and return the access token and refresh token
func (u *UserController) SignUp(c *gin.Context) {
	// Retrieve the body
	signUpUserForm := validator2.SignUpUserForm{}
	if err := c.ShouldBindJSON(&signUpUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(signUpUserForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the tokens
	expiresAt := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
	}

	claimsRefreshToken := &jwt.StandardClaims{
		// Refresh tokens don't expire yet
		Id: guuid.New().String(),
	}

	// Generates access accessToken and refresh accessToken
	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := unSignedToken.SignedString(jwtKey)
	unsignedRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefreshToken)
	refreshToken, errRefreshToken := unsignedRefreshToken.SignedString(jwtKey)

	// Send an error if the tokens didn't sign well
	if err != nil || errRefreshToken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to generate a accessToken",
			"code":  codeErrorServer,
		})
		return
	}

	// Encrypt the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword(bcryptSecret, bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to hash this password",
			"code":  codeErrorServer,
		})
	}

	// Add the user in the database
	user := model.User{Email: signUpUserForm.Email, Password: string(hashedPassword), Name: signUpUserForm.Name}
	database.DB.NewRecord(user)
	err = database.DB.Create(&user).Error

	// Check if there is not an error during the database query
	if err != nil {
		if err.(*pq.Error).Code == "23505" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "The email already exists",
				"code":  codeErrorEmailAlreadyExists,
			})

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while creating the user",
			"code":  codeErrorServer,
		})

		return
	}

	// Save the refresh accessToken
	token := model.Token{Token: refreshToken, DoExpire: false, UserId: user.ID}
	database.DB.NewRecord(token)
	err = database.DB.Create(&token).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error while saving the refresh token",
			"code":  codeErrorServer,
		})

		return
	}

	// Send the tokens
	c.JSONP(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expiresIn":    300000,
	})
}
