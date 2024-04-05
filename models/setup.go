// models/setup.go

package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectDatabase() {
	DB_USERNAME := os.Getenv("DB_USERNAME")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_ENDPOINT := os.Getenv("DB_ENDPOINT")
	DB_NAME := os.Getenv("DB_NAME")

	DB_URL := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USERNAME, DB_PASSWORD, DB_ENDPOINT, DB_NAME)
	database, err := gorm.Open(mysql.Open(DB_URL), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database: " + DB_URL)
	}

	err = database.AutoMigrate(&Book{}, &User{}, &Friendship{}, &Post{}, &PostImages{})

	if err != nil {
		return
	}

	DB = database
}
