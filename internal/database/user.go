package database

import (
	"errors"
	"fmt"
	"log"
	"obs/internal/models"

	"gorm.io/gorm"
)

// CreateUser inserts a new user into the database
func (s *service) CreateUser(user *models.User) error {
	if user == nil {
		log.Println("[DATABASE] CreateUser: received nil user")
		return errors.New("invalid user data")
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		log.Printf("[DATABASE] Error creating user: %v", result.Error)
		return result.Error
	}
	fmt.Println("result: ", result)
	log.Println("[DATABASE] User created successfully")
	return nil
}

// GetUser retrieves a user by ID
func (s *service) GetUser(id uint) (*models.User, error) {
	var user models.User
	result := s.DB.
		Preload("Followers").
		Preload("Following").
		First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("[DATABASE] User not found with ID: %d", id)
			return nil, nil
		}
		log.Printf("[DATABASE] Error retrieving user: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// GetUsers fetches all users
func (s *service) GetUsers() ([]models.User, error) {
	var users []models.User
	result := s.DB.
		Preload("Followers").
		Preload("Following").
		Find(&users)

	if result.Error != nil {
		log.Printf("[DATABASE] Error retrieving users: %v", result.Error)
		return nil, result.Error
	}
	return users, nil
}

// DeleteUser deletes a user by ID
func (s *service) DeleteUser(id uint) error {
	// Permanently delete the user
	result := s.DB.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		log.Printf("[DATABASE] Error deleting user with ID %d: %v", id, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("[DATABASE] No user found with ID %d to delete", id)
		return errors.New("user not found")
	}

	log.Printf("[DATABASE] User with ID %d permanently deleted successfully", id)
	return nil
}

// UpdateUser updates an existing user
func (s *service) UpdateUser(user *models.User) error {
	if user == nil {
		log.Println("[DATABASE] UpdateUser: received nil user")
		return errors.New("invalid user data")
	}

	// Fetch existing user
	existingUser, err := s.GetUser(user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Only update fields that are not zero values
	result := s.DB.Model(&existingUser).Updates(user)
	if result.Error != nil {
		log.Printf("[DATABASE] Error updating user: %v", result.Error)
		return result.Error
	}

	log.Printf("[DATABASE] User ID %d updated successfully", user.ID)
	return nil
}

func (s *service) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := s.DB.Unscoped().Where("LOWER(email) = LOWER(?)", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// FollowUser allows a user to follow another user
func (s *service) FollowUser(followerID, followedID uint) error {
	if followerID == followedID {
		return errors.New("a user cannot follow themselves")
	}

	// Check if the follow relationship already exists
	var existingFollow models.Follow
	result := s.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&existingFollow)
	if result.Error == nil {
		return errors.New("already following this user")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("[DATABASE] Error checking follow relationship: %v", result.Error)
		return result.Error
	}

	// Create a new follow relationship
	follow := models.Follow{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	if err := s.DB.Create(&follow).Error; err != nil {
		log.Printf("[DATABASE] Error following user: %v", err)
		return err
	}

	log.Printf("[DATABASE] User %d followed user %d", followerID, followedID)
	return nil
}

// UnfollowUser allows a user to unfollow another user
func (s *service) UnfollowUser(followerID, followedID uint) error {
	result := s.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&models.Follow{})

	if result.Error != nil {
		log.Printf("[DATABASE] Error unfollowing user: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("not following this user")
	}

	log.Printf("[DATABASE] User %d unfollowed user %d", followerID, followedID)
	return nil
}

// IsFollowing checks if a user is already following another user
func (s *service) IsFollowing(followerID, followedID uint) (bool, error) {
	var follow models.Follow
	err := s.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
