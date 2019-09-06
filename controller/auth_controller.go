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
const codeErrorEmailOrPasswordIncorrect = "CODE_ERROR_EMAIL_OR_PASSWORD_INCORRECT"

var jwtKey = []byte("AllYourBase")

type UserClaim struct {
	User model.User
	jwt.StandardClaims
}

type AuthController struct {
}

// Create the access token with the user information
func createAccessToken(user model.User) (string, error) {
	user.Password = ""
	expiresAt := time.Now().Add(15 * time.Minute)
	claims := UserClaim{
		user,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	// Generates access accessToken and refresh accessToken
	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return unSignedToken.SignedString(jwtKey)
}

// Sign up the user and return the access token and refresh token
func (a *AuthController) SignUp(c *gin.Context) {
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

	// Encrypt the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpUserForm.Password), bcrypt.DefaultCost)
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

	// Create the tokens
	accessToken, err := createAccessToken(user)
	refreshToken, errRefreshToken := guuid.NewUUID()

	// Send an error if the tokens didn't sign well
	if err != nil || errRefreshToken != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to generate an accessToken",
			"code":  codeErrorServer,
		})
		return
	}

	// Save the refresh accessToken
	token := model.Token{Token: refreshToken.String(), UserId: user.ID, IsValid: true}
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
		"expiresIn":    900000,
	})
}

// Login the user and send back the access token and the refresh token
func (a *AuthController) Login(c *gin.Context) {
	// Retrieve the body
	loginUserForm := validator2.LoginUserForm{}
	if err := c.ShouldBindJSON(&loginUserForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(loginUserForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user
	user := model.User{}
	err = database.DB.Where("email = ?", loginUserForm.Email).First(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or password incorrect",
			"code":  codeErrorEmailOrPasswordIncorrect,
		})
		return
	}

	// Check if the password match
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserForm.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email or password incorrect",
			"code":  codeErrorEmailOrPasswordIncorrect,
		})
		return
	}

	// Create the tokens
	accessToken, err := createAccessToken(user)

	// Send an error if the tokens didn't sign well
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to generate an accessToken",
			"code":  codeErrorServer,
		})
		return
	}

	// Retrieve the refresh token
	token := model.Token{}
	err = database.DB.Where("user_id = ?", user.ID).First(&token).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to retrieve the refresh token",
			"code":  codeErrorServer,
		})
		return
	}

	// Send the tokens
	c.JSONP(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": token.Token,
		"expiresIn":    900000,
	})
}

// Refresh the access token thanks to a refresh token
func (a *AuthController) RefreshAccessToken(c *gin.Context) {
	// Retrieve the body
	refreshingTokenForm := validator2.RefreshingTokenForm{}
	if err := c.ShouldBindJSON(&refreshingTokenForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err := validate.Struct(refreshingTokenForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the user linked to the token
	user := model.User{}
	err = database.DB.
		Joins("left join tokens on tokens.user_id = users.id").
		Where("tokens.token = ?", refreshingTokenForm.RefreshToken).
		First(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to retrieve the user",
			"code":  codeErrorServer,
		})
		return
	}

	// Create the access token
	accessToken, err := createAccessToken(user)

	// Send an error if the tokens didn't sign well
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible to generate the access token",
			"code":  codeErrorServer,
		})
		return
	}

	// Send the tokens
	c.JSONP(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"expiresIn":   900000,
	})
}
