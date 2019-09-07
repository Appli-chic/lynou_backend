package service

import (
	"github.com/applichic/lynou/model"
	"github.com/applichic/lynou/util"
)

type UserService struct {
}

// Fetch a user from it's ID
func (u *UserService) FetchUserById(userId interface{}) (model.User, error) {
	user := model.User{}
	err := util.DB.Select("id, email, name").Where("id = ?", userId).First(&user).Error
	return user, err
}

// Fetch a user from it's email
func (u *UserService) FetchUserByEmail(email string) (model.User, error) {
	user := model.User{}
	err := util.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// Fetch a user from the refresh token linked to this account
func (u *UserService) FetchUserFromRefreshToken(refreshToken string) (model.User, error) {
	user := model.User{}
	err := util.DB.
		Joins("left join tokens on tokens.user_id = users.id").
		Where("tokens.token = ?", refreshToken).
		First(&user).Error
	return user, err
}

// Save a user
func (u *UserService) Save(user model.User) (model.User, error) {
	util.DB.NewRecord(user)
	err := util.DB.Create(&user).Error
	return user, err
}
