package database

import (
	"errors"
	"obs/internal/models"

	"gorm.io/gorm"
)

// GetBlogs retrieves all blogs from the database
func (s *service) GetBlogs() ([]models.Blog, error) {
	var blogs []models.Blog
	if err := s.DB.Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// GetBlog fetches a single blog by its ID
func (s *service) GetBlog(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil when the blog is not found
		}
		return nil, err
	}
	return &blog, nil
}

// CreateBlog inserts a new blog into the database
func (s *service) CreateBlog(blog *models.Blog) error {
	return s.DB.Create(blog).Error
}

// DeleteBlog removes a blog by ID after checking its existence
func (s *service) DeleteBlog(id uint) error {
	result := s.DB.Delete(&models.Blog{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return specific error when no rows are deleted
	}
	return nil
}

// UpdateBlog modifies an existing blog's fields
func (s *service) UpdateBlog(blog *models.Blog) error {
	result := s.DB.Model(&models.Blog{}).Where("id = ?", blog.ID).Updates(blog)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return specific error when no rows are updated
	}
	return nil
}
