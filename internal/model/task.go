package model

import (
	"time"

	_ "gorm.io/gorm"
)

// gorm.Model definition
type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

// Task represents a todo task
// @Description Task model represents a todo item in the system
type Task struct {
	Model
	// UserID is the owner of the task
	UserID uint `gorm:"not null" json:"user_id" example:"1"`
	// Event is the main task description
	Event string `gorm:"type:varchar(256);not null" json:"event" example:"Buy groceries"`
	// Completed indicates if the task is done
	Completed bool `gorm:"default:false" json:"completed" example:"false"`
	// IsCycle indicates if the task is recurring
	IsCycle bool `gorm:"default:false" json:"is_cycle" example:"false"`
	// Description provides additional details about the task
	Description string `gorm:"type:text" json:"description" example:"Milk, eggs, bread"`
	// ImportanceLevel indicates task priority (0-5)
	ImportanceLevel int `gorm:"default:0" json:"importance_level" example:"3"`
	// CompletedDate records when the task was completed
	CompletedDate time.Time `gorm:"column:completed_date;default:null" json:"completed_date,omitempty" example:"2025-01-10 15:04:05"`
}
