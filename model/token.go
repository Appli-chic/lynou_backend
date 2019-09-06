package model

import "github.com/jinzhu/gorm"

type Token struct {
	gorm.Model
	Token   string `gorm:"type:varchar(36);unique_index"`
	IsValid bool   `gorm:"not null"`
	UserId  uint   `gorm:"not null"`
}
