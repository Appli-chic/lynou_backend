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
