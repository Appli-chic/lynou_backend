package service

import (
	"github.com/applichic/lynou/model"
	"github.com/applichic/lynou/util"
)

type TokenService struct {
}

// Save the token
func (t *TokenService) Save(token model.Token) (model.Token, error) {
	util.DB.NewRecord(token)
	err := util.DB.Create(&token).Error
	return token, err
}

// Fetch a token the user's email
func (t *TokenService) FetchTokenByUserId(userId interface{}) (model.Token, error) {
	token := model.Token{}
	err := util.DB.Where("user_id = ?", userId).First(&token).Error
	return token, err
}
