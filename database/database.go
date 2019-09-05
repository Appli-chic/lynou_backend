package database

import (
	"github.com/applichic/lynou/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	//conf := config.Get()
	//db, err := gorm.Open("mysql", conf.DSN)
	//
	//if err == nil {
	//	db.DB().SetMaxIdleConns(conf.MaxIdleConn)
	//	DB = db
	//	db.AutoMigrate(&models.AdminUser{})
	//	return db, err
	//}

	//return nil, err
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=lazyos dbname=lynou password=mym2yr sslmode=disable")

	// Send the error
	if err != nil {
		return nil, err
	}

	// Set the database and migrate the models
	db.DB().SetMaxIdleConns(100)
	DB = db
	db.AutoMigrate(&model.User{}, &model.Token{})
	db.Model(&model.Token{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	return db, nil
}
