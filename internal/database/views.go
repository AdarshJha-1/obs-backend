package database

import (
	"errors"
	"obs/internal/models"
)

func (s *service) UpdateView(blogId, userId uint) error {
	if blogId == 0 || userId == 0 {
		return errors.New("invalid blog ID or user ID")
	}

	view := models.View{
		UserID: userId,
		BlogID: blogId,
	}

	// Use `FirstOrCreate` to ensure a unique view per user per blog
	result := s.DB.FirstOrCreate(&view, models.View{UserID: userId, BlogID: blogId})
	return result.Error
}
