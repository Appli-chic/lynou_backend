package util

import (
	"fmt"
	"github.com/applichic/lynou/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	//return nil, err
	dbArgs := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		Conf.DatabaseHost, Conf.DatabasePort, Conf.DatabaseUser, Conf.DatabaseName, Conf.DatabasePassword)
	db, err := gorm.Open(Conf.DatabaseDialect, dbArgs)

	// Send the error
	if err != nil {
		return nil, err
	}

	// Set the database and migrate the models
	db.DB().SetMaxIdleConns(Conf.DatabaseMaxConnection)
	DB = db
	db.AutoMigrate(&model.User{}, &model.Token{})
	db.Model(&model.Token{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	return db, nil
}
