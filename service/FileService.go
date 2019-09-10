package service

import (
	"github.com/applichic/lynou/config"
	"github.com/applichic/lynou/model"
)

type FileService struct {
}

func (f *FileService) Save(file *model.File) error {
	config.DB.NewRecord(file)
	err := config.DB.Create(&file).Error
	return err
}
