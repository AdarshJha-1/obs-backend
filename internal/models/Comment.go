package models

import (
	"time"
)

// Comment model with validation
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content" validate:"required,min=3"`
	UserID    uint      `gorm:"not null;index" json:"user_id" validate:"required"`
	BlogID    uint      `gorm:"not null;index" json:"blog_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	Blog Blog `gorm:"foreignKey:BlogID" json:"-"`
}
