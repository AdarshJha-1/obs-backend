package server

import (
	"errors"
	"log"
	"net/http"
	"obs/internal/models"
	"obs/internal/types"
	"obs/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (s *Server) GetUser(c *gin.Context) {
	strId := c.Param("user_id")
	id64, err := strconv.ParseUint(strId, 10, 64)
	id := uint(id64)
	users, err := s.db.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (s *Server) GetUsers(c *gin.Context) {
	users, err := s.db.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (s *Server) RegisterUser(c *gin.Context) {
	var user models.User

	// Bind JSON request to user model
	if err := c.ShouldBindJSON(&user); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Validate user fields
	if err := user.ValidateUser(); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Validation failed", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Check if user already exists
	existingUser, err := s.db.GetUserByEmail(user.Email)
	if err != nil {
		log.Printf("[DATABASE] Error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, types.Response{
			StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error(),
		})
		return
	}

	if existingUser != nil {
		log.Println("[DATABASE] User already exists with email:", user.Email)
		c.JSON(http.StatusConflict, types.Response{
			StatusCode: http.StatusConflict, Success: false, Message: "Username/Email already taken",
		})
		return
	}
	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Internal server error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	user.Password = string(hashedPassword)

	// Save user to database
	if err := s.db.CreateUser(&user); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error creating user", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	sanitizedUser := struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	// Success response
	res := types.Response{
		StatusCode: http.StatusCreated,
		Success:    true,
		Message:    "User signed up successfully",
		Data:       map[string]any{"user": sanitizedUser},
	}
	c.JSON(http.StatusCreated, res)
}

func (s *Server) Login(c *gin.Context) {
	var login types.SignInModel

	// Bind JSON request
	if err := c.ShouldBindJSON(&login); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Fetch user by email or username
	user, err := s.db.GetUserByEmail(login.Identifier)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Internal server error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Check if user exists
	if user == nil {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Invalid credentials"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Verify password
	if !utils.CheckPassword(login.Password, user.Password) {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Invalid credentials"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Generate JWT token with updated function
	token, err := utils.CreateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error generating token", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Set secure, HTTP-only cookie with JWT token
	c.SetCookie("auth_token", token, 3600, "/", "", true, true)
	// Parameters:
	// - "auth_token": Cookie name
	// - token: The JWT token
	// - 3600: Expiry time (1 hour)
	// - "/": Path (accessible from all routes)
	// - "": Domain (use empty for same-origin or set a specific domain)
	// - true: Secure (send over HTTPS only)
	// - true: HttpOnly (prevent JavaScript access)

	sanitizedUser := struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	// Success response
	res := types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Login successful",
		Data:       map[string]any{"user": sanitizedUser},
	}
	c.JSON(http.StatusOK, res)
}
func (s *Server) DeleteUserById(c *gin.Context) {
	strId := c.Param("user_id")
	id64, err := strconv.ParseUint(strId, 10, 64)
	id := uint(id64)

	if err = s.db.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "user deleted successfully"})
}
func (s *Server) UpdateUserById(c *gin.Context) {
	strId := c.Param("user_id")
	id64, err := strconv.ParseUint(strId, 10, 64)
	id := uint(id64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user models.User

	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user.ID = id
	if err = s.db.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "user updated successfully"})
}

// Blog
// GetAllBlogs handles retrieving all blogs
func (s *Server) GetAllBlogs(c *gin.Context) {
	blogs, err := s.db.GetBlogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch blogs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

// GetBlogByID handles retrieving a blog by ID
func (s *Server) GetBlogByID(c *gin.Context) {
	id, err := parseUintParam(c, "blog_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	blog, err := s.db.GetBlog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching blog"})
		return
	}
	if blog == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"blog": blog})
}

// CreateNewBlog handles creating a new blog
func (s *Server) CreateNewBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog data"})
		return
	}

	if err := s.db.CreateBlog(&blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": "Blog created successfully", "blog": blog})
}

// DeleteBlogByID handles deleting a blog by ID
func (s *Server) DeleteBlogByID(c *gin.Context) {
	id, err := parseUintParam(c, "blog_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	err = s.db.DeleteBlog(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete blog"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// UpdateBlog handles updating a blog
func (s *Server) UpdateBlog(c *gin.Context) {
	id, err := parseUintParam(c, "blog_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog ID"})
		return
	}

	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid blog data"})
		return
	}

	blog.ID = id // Assign the extracted ID

	err = s.db.UpdateBlog(&blog)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Blog updated successfully", "blog": blog})
}

// parseUintParam extracts and converts a URL parameter to uint
func parseUintParam(c *gin.Context, param string) (uint, error) {
	strID := c.Param(param)
	id64, err := strconv.ParseUint(strID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}
