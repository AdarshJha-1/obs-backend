package database

import (
	"errors"
	"obs/internal/models"

	"gorm.io/gorm"
)

// Get the total number of likes for a blog
func (s *service) GetLikesForBlog(blogID uint) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Like{}).Where("blog_id = ?", blogID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Get a like entry by its ID
func (s *service) GetLikeByID(likeID uint) (*models.Like, error) {
	var like models.Like
	err := s.DB.First(&like, likeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("like not found")
	}
	return &like, err
}

// Like a blog (with duplicate check)
func (s *service) LikeBlog(like *models.Like) error {
	// Check if the blog exists
	var blog models.Blog
	if err := s.DB.First(&blog, like.BlogID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("blog not found")
	}

	// Check if the user has already liked this blog
	var existingLike models.Like
	err := s.DB.Where("user_id = ? AND blog_id = ?", like.UserID, like.BlogID).First(&existingLike).Error
	if err == nil {
		return errors.New("user already liked this blog")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err // Return DB error if any
	}

	// Create a new like entry
	if err := s.DB.Create(like).Error; err != nil {
		return err
	}
	return nil
}

// Unlike a blog (with existence check)
func (s *service) UnlikeBlog(userID, blogID uint) error {
	// Check if the like exists
	var like models.Like
	err := s.DB.Where("user_id = ? AND blog_id = ?", userID, blogID).First(&like).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("like not found")
	} else if err != nil {
		return err
	}

	// Delete the like entry
	if err := s.DB.Delete(&like).Error; err != nil {
		return err
	}
	return nil
}

// Delete a like by its ID (checks if like exists before deleting)
func (s *service) DeleteLike(likeID uint) error {
	var like models.Like
	err := s.DB.First(&like, likeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("like not found")
	} else if err != nil {
		return err
	}

	if err := s.DB.Delete(&like).Error; err != nil {
		return err
	}
	return nil
}

// GetLikeByUserAndBlog retrieves a like entry for a given user and blog
func (s *service) GetLikeByUserAndBlog(userID, blogID uint) (*models.Like, error) {
	var like models.Like
	err := s.DB.Where("user_id = ? AND blog_id = ?", userID, blogID).First(&like).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("like not found")
	}
	return &like, err
}
