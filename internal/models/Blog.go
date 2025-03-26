package models

import "time"

// Blog model
type Blog struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:225;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Views     int       `gorm:"default:0" json:"views"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID" json:"-"`
	Comments []Comment `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE;" json:"comments"`
	Likes    []Like    `gorm:"foreignKey:BlogID;constraint:OnDelete:CASCADE;" json:"likes"`
}
