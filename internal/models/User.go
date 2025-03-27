package models

import (
	"time"
)

// User model with validation
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:100;not null" json:"username" validate:"required,min=3,max=100"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"type:text;not null" json:"password" validate:"required,min=6"`
	Role      string    `gorm:"size:50;not null;default:'author'" json:"role" validate:"required,oneof=author admin"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Blogs     []Blog    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Comments  []Comment `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Likes     []Like    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Followers []Follow  `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE;"`
	Following []Follow  `gorm:"foreignKey:FollowedID;constraint:OnDelete:CASCADE;"`
}
