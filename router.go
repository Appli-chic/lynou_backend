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
	authController := controller.NewAuthController()
	userController := controller.NewUserController()
	postController := controller.NewPostController()
	storageController := new(controller.StorageController)

	api := router.Group("/api")
	{
		// Auth routes
		api.POST("/auth", authController.SignUp)
		api.POST("/auth/login", authController.Login)
		api.POST("/auth/refresh", authController.RefreshAccessToken)

		// Retrieve the videos
		api.GET("/video/:name", storageController.DownloadVideoFile)

		// Need to be logged in routes
		loggedInGroup := api.Group("/")
		loggedInGroup.Use(util.AuthenticationRequired())
		{
			// User
			loggedInGroup.GET("/user", userController.FetchUser)
			loggedInGroup.GET("/user/photo", userController.FetchProfilePhoto)

			// Post
			loggedInGroup.POST("/post", postController.CreatePost)
			loggedInGroup.GET("/posts", postController.FetchPosts)

			// Storage
			loggedInGroup.GET("/file/:name", storageController.DownloadFile)
			loggedInGroup.POST("/file/:name", storageController.UploadFile)
		}
	}

	return router
}
