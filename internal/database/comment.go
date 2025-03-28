package database

import (
	"errors"
	"obs/internal/models"

	"gorm.io/gorm"
)

// GetComments retrieves all comments for a blog with user info
func (s *service) GetComments(blogID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := s.DB.Preload("User").Where("blog_id = ?", blogID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetComment fetches a single comment by its ID
func (s *service) GetComment(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := s.DB.Preload("User").First(&comment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// CreateComment inserts a new comment into the database
func (s *service) CreateComment(comment *models.Comment) error {
	return s.DB.Create(comment).Error
}

// UpdateComment updates only the content of a comment
func (s *service) UpdateComment(id uint, userID uint, content string) error {
	result := s.DB.Model(&models.Comment{}).Where("id = ? AND user_id = ?", id, userID).Update("content", content)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteComment ensures only the comment owner can delete it
func (s *service) DeleteComment(id uint, userID uint) error {
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
