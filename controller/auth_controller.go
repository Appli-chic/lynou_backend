package controller

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
	"github.com/applichic/lynou/service"
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

type UserClaim struct {
	User model.User
	jwt.StandardClaims
}

type AuthController struct {
	userService  *service.UserService
	tokenService *service.TokenService
}

func NewAuthController() *AuthController {
	authController := new(AuthController)
	authController.userService = new(service.UserService)
	authController.tokenService = new(service.TokenService)
	return authController
}

// Create the access token with the service information
func createAccessToken(user model.User) (string, error) {
	var newUser = model.User{}
	newUser.ID = user.ID
	expiresAt := time.Now().Add(time.Duration(config.Conf.JwtTokenExpiration) * time.Millisecond)
	claims := UserClaim{
		newUser,
		jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	// Generates access accessToken and refresh accessToken
	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return unSignedToken.SignedString([]byte(config.Conf.JwtSecret))
}

// Sign up the service and return the access token and refresh token
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

	// Encrypt the service's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpUserForm.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
	}

	// Add the service in the database
	user := model.User{Email: signUpUserForm.Email, Password: string(hashedPassword), Name: signUpUserForm.Name, Photo: "profile.png"}
	err = a.userService.Save(&user)

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
			"error": err.Error(),
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
	token, err = a.tokenService.Save(token)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})

		return
	}

	// Send the tokens
	c.JSONP(http.StatusCreated, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"expiresIn":    config.Conf.JwtTokenExpiration,
	})
}

// Login the service and send back the access token and the refresh token
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

	// Find the service
	user, err := a.userService.FetchUserByEmail(loginUserForm.Email)

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
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Retrieve the refresh token
	token, err := a.tokenService.FetchTokenByUserId(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the tokens
	c.JSONP(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": token.Token,
		"expiresIn":    config.Conf.JwtTokenExpiration,
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

	// Get the service linked to the token
	user, err := a.userService.FetchUserFromRefreshToken(refreshingTokenForm.RefreshToken)
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
		"expiresIn":   config.Conf.JwtTokenExpiration,
	})
}
