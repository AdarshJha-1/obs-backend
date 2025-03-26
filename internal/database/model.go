package database

import (
	"time"
)

// User model
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"size:100;unique;not null"`
	Email     string `gorm:"size:100;unique;not null"`
	Password  string `gorm:"type:text;not null"` // Store hashed passwords
	Role      string `gorm:"size:50;not null;default:'author'"`
	CreatedAt time.Time

	// Relationships
	Blogs     []Blog    `gorm:"foreignKey:UserID"`
	Comments  []Comment `gorm:"foreignKey:UserID"`
	Likes     []Like    `gorm:"foreignKey:UserID"`
	Followers []Follow  `gorm:"foreignKey:FollowerID"`
	Following []Follow  `gorm:"foreignKey:FollowedID"`
}

// Blog model
type Blog struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:225;not null"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint   `gorm:"not null;index"` // Foreign key reference to User
	Views     int    `gorm:"default:0"`
	CreatedAt time.Time

	// Relationships
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:BlogID"`
	Likes    []Like    `gorm:"foreignKey:BlogID"`
}

// Comment model
type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `gorm:"type:text;not null"`
	UserID    uint   `gorm:"not null;index"`
	BlogID    uint   `gorm:"not null;index"`
	CreatedAt time.Time

	// Relationships
	User User `gorm:"foreignKey:UserID"`
	Blog Blog `gorm:"foreignKey:BlogID"`
}

// Like model
type Like struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null;index"`
	BlogID    uint `gorm:"not null;index"`
	CreatedAt time.Time

	// Relationships
	User User `gorm:"foreignKey:UserID"`
	Blog Blog `gorm:"foreignKey:BlogID"`
}

// Follow model
type Follow struct {
	ID         uint `gorm:"primaryKey"`
	FollowerID uint `gorm:"not null;index"` // User who follows
	FollowedID uint `gorm:"not null;index"` // User who is followed
	CreatedAt  time.Time

	// Relationships
	Follower User `gorm:"foreignKey:FollowerID"`
	Followed User `gorm:"foreignKey:FollowedID"`
}
