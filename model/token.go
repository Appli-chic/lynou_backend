package model

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	Token    string `gorm:"type:varchar(255);unique_index"`
	DoExpire bool
	UserId   uint
}
