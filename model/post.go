package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Text   string `gorm:"type:text;not null"`
	UserId uint   `gorm:"not null"`
}
