package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConn() *gorm.DB {

	dsn := "masih:borotosh@tcp(127.0.0.1:3306)/hotelbooking?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	return db

}
