package utils

import "obs/internal/models"

type SanitizedUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Pfp       string `json:"pfp"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

func SanitizedUserData(user *models.User) SanitizedUser {
	return SanitizedUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Pfp:       user.Pfp,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
