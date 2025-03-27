package database

import (
	"errors"
	"obs/internal/models"

	"gorm.io/gorm"
)

// Get like count for a blog
func (s *service) GetLikesForBlog(blogID uint) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Like{}).Where("blog_id = ?", blogID).Count(&count).Error
	return count, err
}

// Get a like by ID
func (s *service) GetLikeByID(likeID uint) (*models.Like, error) {
	var like models.Like
	err := s.DB.First(&like, likeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &like, err
}

// Like a blog
func (s *service) LikeBlog(like *models.Like) error {
	return s.DB.Create(like).Error
}

// Unlike a blog
func (s *service) UnlikeBlog(userID, blogID uint) error {
	return s.DB.Where("user_id = ? AND blog_id = ?", userID, blogID).Delete(&models.Like{}).Error
}

// Delete a like by ID
func (s *service) DeleteLike(likeID uint) error {
	return s.DB.Delete(&models.Like{}, likeID).Error
}
