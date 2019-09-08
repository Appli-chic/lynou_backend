package service

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
)

type TokenService struct {
}

// Save the token
func (t *TokenService) Save(token model.Token) (model.Token, error) {
	config.DB.NewRecord(token)
	err := config.DB.Create(&token).Error
	return token, err
}

// Fetch a token the user's email
func (t *TokenService) FetchTokenByUserId(userId interface{}) (model.Token, error) {
	token := model.Token{}
	err := config.DB.Where("user_id = ?", userId).First(&token).Error
	return token, err
}
