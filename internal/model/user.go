package model

import "gorm.io/gorm"

// User 用户模型
// @Description 用户信息模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(32);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null"`
}
