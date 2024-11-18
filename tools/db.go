package tools

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"vidspark/configs"
)

func InitDB() *gorm.DB {
	config := configs.InitConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MySql.User, config.MySql.Password, config.MySql.Host, config.MySql.Port, config.MySql.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	return db
}
