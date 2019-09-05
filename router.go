package main

import (
	"github.com/applichic/lynou/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		userController := new(controller.UserController)
		api.POST("/auth", userController.SignUp)
	}

	return router
}
