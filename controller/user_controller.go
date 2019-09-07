package controller

import (
	"github.com/applichic/lynou/service"
	"github.com/applichic/lynou/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	userController := new(UserController)
	userController.userService = new(service.UserService)
	return userController
}

// Fetch service's data
func (u *UserController) FetchUser(c *gin.Context) {
	token, _ := util.GetToken(c)
	userClaims := token.Claims.(jwt.MapClaims)["User"].(map[string]interface{})
	user, err := u.userService.FetchUserById(userClaims["ID"])

	// Check if the service exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the service information
	c.JSON(http.StatusOK, gin.H{
		"service": user,
	})
}
