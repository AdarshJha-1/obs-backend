package models

import (
	"time"
)

// Follow model with validation
type Follow struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FollowerID uint      `gorm:"not null;index" json:"follower_id" validate:"required,nefield=FollowedID"`
	FollowedID uint      `gorm:"not null;index" json:"followed_id" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Follower User `gorm:"foreignKey:FollowerID" json:"-"`
	Followed User `gorm:"foreignKey:FollowedID" json:"-"`
}
