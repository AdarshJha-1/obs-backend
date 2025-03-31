package database

import (
	"errors"
	"obs/internal/models"

	"gorm.io/gorm"
)

// GetBlogs retrieves all blogs along with their related data
func (s *service) GetBlogs() ([]models.Blog, error) {
	var blogs []models.Blog
	if err := s.DB.Preload("User").Preload("Comments").Preload("Likes").Preload("Views").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// GetBlog fetches a single blog by its ID along with its related data
func (s *service) GetBlog(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.Preload("User").Preload("Comments").Preload("Likes").First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil when the blog is not found
		}
		return nil, err
	}
	return &blog, nil
}

// CreateBlog inserts a new blog into the database and returns it
func (s *service) CreateBlog(blog *models.Blog) (*models.Blog, error) {
	if err := s.DB.Create(blog).Error; err != nil {
		return nil, err
	}
	return blog, nil
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

// UpdateBlog modifies an existing blog's fields safely
func (s *service) UpdateBlog(blog *models.Blog) error {
	result := s.DB.Model(&models.Blog{}).Where("id = ?", blog.ID).Updates(map[string]interface{}{
		"title":   blog.Title,
		"content": blog.Content,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Return specific error when no rows are updated
	}
	return nil
}
