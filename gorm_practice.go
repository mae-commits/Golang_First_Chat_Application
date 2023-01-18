package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

func main() {
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}
	db.AutoMigrate(&User{})
	user := User{}
	var count int64
	db.Model(user).Where("name = ?", "username").Count(&count)
	fmt.Println(count)
}
