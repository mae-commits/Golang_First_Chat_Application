package account

import (
	"chatapp/domain"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ユーザ認証
func GetUser(userName string, password string) (count int64) {
	db, err := gorm.Open(sqlite.Open("userData.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database.")
	}
	db.AutoMigrate(&domain.User{})
	db.Model(&domain.User{}).Where("name = ?", userName).Where("password = ?", password).Count(&count)
	return count
}
