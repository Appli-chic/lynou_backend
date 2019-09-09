package controller

import (
	"github.com/applichic/lynou/model"
	"github.com/applichic/lynou/service"
	"github.com/applichic/lynou/util"
	validator2 "github.com/applichic/lynou/validator"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

type PostController struct {
	postService *service.PostService
	userService *service.UserService
}

func NewPostController() *PostController {
	PostController := new(PostController)
	PostController.postService = new(service.PostService)
	PostController.userService = new(service.UserService)
	return PostController
}

func (p *PostController) FetchPosts(c *gin.Context) {
	user, err := util.GetUserFromToken(c)

	// Check if the user exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Retrieve the page argument
	pageString := c.DefaultQuery("page", "0")
	page, err := strconv.Atoi(pageString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Retrieve the posts
	posts, err := p.postService.FetchWallPosts(user.ID, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the posts
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

// Create a new post
func (p *PostController) CreatePost(c *gin.Context) {
	user, err := util.GetUserFromToken(c)

	// Check if the user exists
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Retrieve the body
	postForm := validator2.PostForm{}
	if err := c.ShouldBindJSON(&postForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the form
	validate := validator.New()
	err = validate.Struct(postForm)

	// Check if the form is valid
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the post
	post := model.Post{Text: postForm.Text, UserId: user.ID}
	post, err = p.postService.Save(post)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  codeErrorServer,
		})
		return
	}

	// Send the post
	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}
