package models

import (
	"time"
)

// Blog model with validation
type Blog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:225;not null" json:"title" validate:"required,min=3,max=225"`
	Content   string    `gorm:"type:text;not null" json:"content" validate:"required,min=10"`
	UserID    uint      `gorm:"not null;index" json:"user_id" validate:"required"`
	Views     int       `gorm:"default:0" json:"views" validate:"gte=0"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Comments []Comment `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE;" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE;" json:"likes"`
}
