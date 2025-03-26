package database

import "obs/internal/models"

func (s *service) GetBlogs() ([]models.Blog, error) {
	var blogs []models.Blog
	result := s.DB.Find(&blogs)
	if result.Error != nil {
		return nil, result.Error
	}
	return blogs, nil
}

func (s *service) CreateBlog(blog models.Blog) error {
	result := s.DB.Create(&blog)
	return result.Error
}

func (s *service) GetBlog(id uint64) (*models.Blog, error) {
	var blog models.Blog
	result := s.DB.First(&blog, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &blog, nil
}

func (s *service) DeleteBlog(id uint64) error {
	result := s.DB.Delete(&models.Blog{}, id)
	return result.Error
}

func (s *service) UpdateBlog(blog models.Blog) error {
	result := s.DB.Save(&blog)
	return result.Error
}
