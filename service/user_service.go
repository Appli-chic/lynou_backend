package service

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
)

type UserService struct {
}

// Fetch a user from it's ID
func (u *UserService) FetchUserById(userId interface{}) (model.User, error) {
	user := model.User{}
	err := config.DB.Select("id, email, name, photo").Where("id = ?", userId).First(&user).Error
	return user, err
}

// Fetch a user from it's email
func (u *UserService) FetchUserByEmail(email string) (model.User, error) {
	user := model.User{}
	err := config.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// Fetch a user from the refresh token linked to this account
func (u *UserService) FetchUserFromRefreshToken(refreshToken string) (model.User, error) {
	user := model.User{}
	err := config.DB.
		Joins("left join tokens on tokens.user_id = users.id").
		Where("tokens.token = ?", refreshToken).
		First(&user).Error
	return user, err
}

// Save a user
func (u *UserService) Save(user *model.User) error {
	config.DB.NewRecord(user)
	err := config.DB.Create(&user).Error
	return err
}

// Fetch the user's photo with the user's id
func (u *UserService) FetchProfilePhotoByUserId(userId interface{}) (model.User, error) {
	user := model.User{}
	err := config.DB.Select("photo").Where("id = ?", userId).First(&user).Error
	return user, err
}
