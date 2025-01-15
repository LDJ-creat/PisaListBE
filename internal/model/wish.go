package model

import (
	"time"
)

// Wish 心愿模型
// @Description 用户的心愿信息
type Wish struct {
	ID          uint       `json:"id" gorm:"primarykey" example:"1"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-01-10T15:04:05Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-01-10T15:04:05Z"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" example:"2024-01-10T15:04:05Z"`
	UserID      uint       `json:"user_id" gorm:"not null" example:"1"`
	Event       string     `json:"event" gorm:"type:varchar(256);not null" example:"环游世界"`
	IsCycle     bool       `json:"is_cycle" gorm:"default:false" example:"false"`
	Description string     `json:"description" gorm:"type:text" example:"想去看看世界的每个角落"`
	IsShared    bool       `json:"is_shared" gorm:"default:false" example:"false"`
}

// SharedWish 共享心愿模型
// @Description 用户分享到社区的心愿信息
type SharedWish struct {
	ID             uint       `json:"id" gorm:"primarykey" example:"1"`
	CreatedAt      time.Time  `json:"created_at" example:"2024-01-10T15:04:05Z"`
	UpdatedAt      time.Time  `json:"updated_at" example:"2024-01-10T15:04:05Z"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" example:"2024-01-10T15:04:05Z"`
	OriginalWishID uint       `json:"original_wish_id" gorm:"not null" example:"1"`
	Event          string     `json:"event" gorm:"type:varchar(256);not null" example:"环游世界"`
	Description    string     `json:"description" gorm:"type:text" example:"想去看看世界的每个角落"`
	SharedByUserID uint       `json:"shared_by_user_id" gorm:"not null" example:"1"`
}
