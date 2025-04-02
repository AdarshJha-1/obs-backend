package database

import (
	"errors"
	"gorm.io/gorm"
	"obs/internal/models"
)

// AdminGetUsers retrieves all users from the database
func (s *service) AdminGetUsers() ([]models.User, error) {
	var users []models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// AdminGetUser retrieves a single user by their ID
func (s *service) AdminGetUser(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// AdminDeleteUser deletes a user by ID
func (s *service) AdminDeleteUser(id uint) error {
	result := s.DB.Unscoped().Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AdminUpdateUser updates an existing user's information
func (s *service) AdminUpdateUser(user *models.User) error {
	result := s.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AdminGetBlogs retrieves all blogs along with their related data
func (s *service) AdminGetBlogs() ([]models.Blog, error) {
	var blogs []models.Blog
	if err := s.DB.Preload("Author").Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}

// AdminGetBlog retrieves a single blog by its ID along with related data
func (s *service) AdminGetBlog(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := s.DB.Preload("Author").First(&blog, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &blog, nil
}

// AdminDeleteBlog deletes a blog by ID
func (s *service) AdminDeleteBlog(id uint) error {
	result := s.DB.Delete(&models.Blog{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AdminUpdateBlog updates a blog's information
func (s *service) AdminUpdateBlog(blog *models.Blog) error {
	result := s.DB.Model(&models.Blog{}).Where("id = ?", blog.ID).Updates(map[string]interface{}{
		"title":   blog.Title,
		"content": blog.Content,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AdminDeleteComment deletes a comment by ID
func (s *service) AdminDeleteComment(id uint) error {
	result := s.DB.Delete(&models.Comment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// AdminUpdateComment updates a comment's content
func (s *service) AdminUpdateComment(id uint, content string) error {
	result := s.DB.Model(&models.Comment{}).Where("id = ?", id).Update("content", content)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (s *service) AdminGetComments() ([]models.Comment, error) {
	var comments []models.Comment
	if err := s.DB.Find(&comments).Error; err != nil {
		return nil, err // Return the error if the query fails
	}
	return comments, nil
}

type DashboardData struct {
	TotalUsers    int64 `json:"total_users"`
	TotalBlogs    int64 `json:"total_blogs"`
	TotalComments int64 `json:"total_comments"`
	// Add other metrics as needed
}

// GetAdminDashboardData retrieves data for the admin dashboard (total users, blogs, comments)
func (s *service) GetAdminDashboardData() (DashboardData, error) {
	var dashboardData DashboardData

	// Query to get the total number of users
	if err := s.DB.Model(&models.User{}).Count(&dashboardData.TotalUsers).Error; err != nil {
		return dashboardData, err
	}

	// Query to get the total number of blogs
	if err := s.DB.Model(&models.Blog{}).Count(&dashboardData.TotalBlogs).Error; err != nil {
		return dashboardData, err
	}

	// Query to get the total number of comments
	if err := s.DB.Model(&models.Comment{}).Count(&dashboardData.TotalComments).Error; err != nil {
		return dashboardData, err
	}

	// Return the dashboard data
	return dashboardData, nil
}
