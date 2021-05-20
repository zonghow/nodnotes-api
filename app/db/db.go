package db

import (
	"fmt"
	"nodnotes-api/app/config"
	"nodnotes-api/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var D = struct {
	Client *gorm.DB
}{}

func init() {
	gormConfig := config.C.Gorm
	db, err := gorm.Open(mysql.Open(gormConfig.Dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("Connect to db failed: %s \n", err))
	}
	D.Client = db
	if gormConfig.AutoMigrate {
		D.Client.AutoMigrate(&models.UserModel{})
		D.Client.AutoMigrate(&models.NodeModel{})
	}
}
