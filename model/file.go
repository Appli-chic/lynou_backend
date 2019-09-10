package model

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Name      string `gorm:"type:varchar(50);not null"`
	Thumbnail string `gorm:"type:varchar(50)"`
	Type      uint   `gorm:"not null"`

	PostId uint  `gorm:"not null"`
	Post   *Post `gorm:"foreignkey:PostId"`
}
