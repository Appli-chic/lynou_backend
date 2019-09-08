package controller

import (
	"github.com/applichic/lynou/service"
	"github.com/applichic/lynou/util"
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

// Fetch profile's photo
func (u *UserController) FetchProfilePhoto(c *gin.Context) {
	user, err := util.GetUserFromToken(c)

	// Check if the user exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the user information
	c.JSON(http.StatusOK, gin.H{
		"photo": user.Photo,
	})
}

// Fetch user's data
func (u *UserController) FetchUser(c *gin.Context) {
	user, err := util.GetUserFromToken(c)

	// Check if the user exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the user information
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
