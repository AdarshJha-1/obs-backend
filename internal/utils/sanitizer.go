package utils

import "obs/internal/models"

type SanitizedUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Pfp       string `json:"pfp"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	Followers []uint `json:"followers"` // List of follower IDs
	Following []uint `json:"following"` // List of following IDs
}

func SanitizedUserData(user *models.User) SanitizedUser {
	// Extract follower and following IDs
	var followers []uint
	var following []uint

	// `Followers` now contains `User` records instead of `Follow`
	for _, follower := range user.Followers {
		followers = append(followers, follower.ID) // Extracting User ID
	}

	// `Following` now contains `User` records instead of `Follow`
	for _, followingUser := range user.Following {
		following = append(following, followingUser.ID) // Extracting User ID
	}

	return SanitizedUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Pfp:       user.Pfp,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		Followers: followers,
		Following: following,
	}
}
