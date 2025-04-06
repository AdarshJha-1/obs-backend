package server

import (
	"log"
	"net/http"
	"obs/internal/models"
	"obs/internal/types"
	"obs/internal/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	existingUser, err := s.db.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("[DATABASE] Error checking user existence: %v", err)
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if existingUser != nil {
		res := types.Response{StatusCode: http.StatusConflict, Success: false, Message: "Username/Email already taken"}
		c.JSON(http.StatusConflict, res)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Internal server error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	user.Password = string(hashedPassword)

	if err := s.db.CreateUser(&user); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error creating user", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusCreated, Success: true, Message: "User signed up successfully", Data: map[string]any{"user": utils.SanitizedUserData(&user)}}
	c.JSON(http.StatusCreated, res)
}

// LoginUser handles user authentication, sets a cookie, and checks if user is an author or admin
func (s *Server) LoginUser(c *gin.Context) {
	var creds types.SignInModel
	if err := c.ShouldBindJSON(&creds); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := s.db.GetUserByEmail(creds.Identifier)
	if err != nil || user == nil {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Invalid credentials"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Check password
	if !utils.CheckPassword(creds.Password, user.Password) {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Invalid credentials"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Check if the user is an admin or author
	var role string
	switch user.Role {
	case "admin":
		role = "admin"
	default:
		role = "author"
	}

	token, err := utils.CreateJWT(user.ID, user.Username, user.Email, role) // pass role to JWT creation
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error generating token", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Set token as cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		SameSite: http.SameSiteLaxMode,
	})

	// Optionally sanitize user data before returning it
	sanitizedUser := utils.SanitizedUserData(user)
	res := types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Login successful",
		Data:       map[string]any{"user": sanitizedUser, "role": role},
	}
	c.JSON(http.StatusOK, res)
}

// LogoutUser handles user logout by clearing the auth_token cookie
func (s *Server) LogoutUser(c *gin.Context) {
	// Clear the authentication cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   -1, // Expire the cookie immediately
		SameSite: http.SameSiteLaxMode,
	})

	res := types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Logout successful",
	}
	c.JSON(http.StatusOK, res)
}
