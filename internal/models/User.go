package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID        uint64         `gorm:"primaryKey"`
	Username  string         `gorm:"size:100;unique;not null" json:"username"`
	Email     string         `gorm:"size:100;unique;not null" json:"email"`
	Password  string         `gorm:"type:text;not null" json:"password"`
	Role      string         `gorm:"size:50;not null;default:'author'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Blogs     []Blog    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Comments  []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Likes     []Like    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Followers []Follow  `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE;"`
	Following []Follow  `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE;"`
}
