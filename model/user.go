package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email    string  `gorm:"type:varchar(255);unique_index"`
	Password string  `gorm:"type:varchar(64);not null"`
	Name     string  `gorm:"type:varchar(100);not null"`
	Photo    string  `gorm:"type:varchar(255);not null"`
	Tokens   []Token `gorm:"foreignkey:UserRefer"`
}
