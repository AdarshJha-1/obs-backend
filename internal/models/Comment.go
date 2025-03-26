package models

import "time"

// Comment model
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	BlogID    uint      `gorm:"not null;index" json:"blog_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
	Blog Blog `gorm:"foreignKey:BlogID" json:"-"`
}
