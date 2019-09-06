package main

import (
	"github.com/applichic/lynou/controller"
	"github.com/applichic/lynou/util"
	"github.com/gin-gonic/gin"
)

// Creates all the routes
func InitRouter() *gin.Engine {
	router := gin.Default()

	// Create controllers
	authController := new(controller.AuthController)
	userController := new(controller.UserController)
	storageController := new(controller.StorageController)

	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/auth", authController.SignUp)
		api.POST("/auth/login", authController.Login)
		api.GET("/auth/refresh", authController.RefreshAccessToken)

		// Need to be logged in routes
		loggedInGroup := api.Group("/")
		loggedInGroup.Use(util.AuthenticationRequired())
		{
			loggedInGroup.GET("/user", userController.FetchUser)
			loggedInGroup.GET("/file/:path", storageController.DownloadImage)
			loggedInGroup.PUT("/file/:name", storageController.UploadFile)
		}
	}

	return router
}
