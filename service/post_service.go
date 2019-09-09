package service

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
)

type PostService struct {
}

// Save a post
func (p *PostService) Save(post model.Post) (model.Post, error) {
	config.DB.NewRecord(post)
	err := config.DB.Create(&post).Error
	return post, err
}

// Fetch wall posts
func (p *PostService) FetchWallPosts(userId interface{}, page int) ([]model.Post, error) {
	nbRows := 10
	var posts []model.Post
	err := config.DB.
		Joins("left join users on users.id = user_id").
		Limit(nbRows).
		Offset(page * nbRows).
		Preload("User").
		Order("created_at desc").
		Find(&posts).Error
	return posts, err
}
