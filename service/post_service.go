package service

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
)

type PostService struct {
}

// Save a post
func (p *PostService) Save(post *model.Post) error {
	config.DB.NewRecord(post)
	err := config.DB.Create(&post).Error
	return err
}

// Fetch wall posts
func (p *PostService) FetchWallPosts(userId interface{}, page int) ([]model.Post, error) {
	nbRows := 10
	var posts []model.Post
	err := config.DB.
		Limit(nbRows).
		Offset(page * nbRows).
		Preload("User").
		Preload("Files").
		Order("created_at desc").
		Find(&posts).Error

	// Remove the hashed password in the user
	for _, post := range posts {
		post.User.Password = ""
	}

	return posts, err
}
