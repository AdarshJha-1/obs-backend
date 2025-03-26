package database

import (
	"log"
	"obs/internal/models"
)

func (s *service) CreateUser(user *models.User) error {
	result := s.DB.Create(user)
	if result.Error != nil {
		log.Printf("[DATABASE] ❌ Error creating user: %v", result.Error)
		return result.Error
	}
	log.Println("[DATABASE] ✅ User created successfully!")
	return nil
}

func (s *service) GetUser(id uint64) (*models.User, error) {
	var user models.User
	result := s.DB.First(&user, id)
	return &user, result.Error
}
func (s *service) GetUsers() ([]models.User, error) {
	var users []models.User
	result := s.DB.Find(&users)
	return users, result.Error
}

func (s *service) DeleteUser(id uint64) error {
	result := s.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (s *service) UpdateUser(user models.User) error {
	result := s.DB.Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
