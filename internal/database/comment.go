package database

import (
	"errors"
	"obs/internal/models"

	"gorm.io/gorm"
)

// GetComments retrieves all comments for a given blog ID
func (s *service) GetComments(blogID uint) ([]models.Comment, error) {
	var comments []models.Comment
	if err := s.DB.Where("blog_id = ?", blogID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// GetComment fetches a single comment by its ID
func (s *service) GetComment(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := s.DB.First(&comment, id).Error; err != nil {
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

// UpdateComment modifies an existing comment
func (s *service) UpdateComment(comment *models.Comment) error {
	result := s.DB.Model(&models.Comment{}).Where("id = ?", comment.ID).Updates(comment)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteComment removes a comment by ID
func (s *service) DeleteComment(id uint) error {
	result := s.DB.Delete(&models.Comment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
