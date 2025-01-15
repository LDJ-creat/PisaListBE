package model

import (
	"time"
)

// User 用户模型
// @Description 用户信息模型
type User struct {
	ID        uint       `json:"id" gorm:"primarykey" example:"1"`
	CreatedAt time.Time  `json:"created_at" example:"2024-01-10T15:04:05Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2024-01-10T15:04:05Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:"2024-01-10T15:04:05Z"`
	Username  string     `json:"username" gorm:"type:varchar(32);uniqueIndex;not null" example:"johndoe"`
	Password  string     `json:"-" gorm:"type:varchar(256);not null"`
	Email     string     `json:"email" gorm:"type:varchar(256);uniqueIndex" example:"john@example.com"`
}
