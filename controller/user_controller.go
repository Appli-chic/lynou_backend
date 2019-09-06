package controller

import (
	"github.com/applichic/lynou/model"
	"github.com/applichic/lynou/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
}

// Fetch user's data
func (u *UserController) FetchUser(c *gin.Context) {
	token, _ := util.GetToken(c)
	userClaims := token.Claims.(jwt.MapClaims)["User"].(map[string]interface{})

	user := model.User{}
	err := util.DB.Select("id, email, name").Where("id = ?", userClaims["ID"]).First(&user).Error

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
