package models

import (
	"time"
)

type View struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:user_blog_unique" json:"user_id"`
	BlogID    uint      `gorm:"not null;index;uniqueIndex:user_blog_unique" json:"blog_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	Blog Blog `gorm:"foreignKey:BlogID" json:"-"`
}
