package models

import "time"

// Follow model
type Follow struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FollowerID uint      `gorm:"not null;index" json:"follower_id"`
	FollowedID uint      `gorm:"not null;index" json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Follower User `gorm:"foreignKey:FollowerID" json:"-"`
	Followed User `gorm:"foreignKey:FollowedID" json:"-"`
}
