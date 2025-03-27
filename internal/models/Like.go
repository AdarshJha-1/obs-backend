package models

import (
	"time"
)

// Like model with validation
type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id" validate:"required"`
	BlogID    uint      `gorm:"not null;index" json:"blog_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	Blog Blog `gorm:"foreignKey:BlogID" json:"-"`
}
