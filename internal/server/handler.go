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

type sanitizedUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

func sanitizedUserData(user *models.User) sanitizedUser {
	return sanitizedUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func (s *Server) GetUser(c *gin.Context) {
	id, err := parseUintParam(c, "user_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid user ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := s.db.GetUser(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	if user == nil {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "User not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "User fetched successfully",
		Data:       map[string]any{"user": sanitizedUserData(user)},
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) GetUsers(c *gin.Context) {
	users, err := s.db.GetUsers()
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	sanitizedUsers := make([]sanitizedUser, len(users))
	for i, user := range users {
		sanitizedUsers[i] = sanitizedUserData(&user)
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Users fetched successfully", Data: map[string]any{"users": sanitizedUsers}}
	c.JSON(http.StatusOK, res)
}

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

	res := types.Response{StatusCode: http.StatusCreated, Success: true, Message: "User signed up successfully", Data: map[string]any{"user": sanitizedUserData(&user)}}
	c.JSON(http.StatusCreated, res)
}

// LoginUser handles user authentication and sets a cookie
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

	if !utils.CheckPassword(creds.Password, user.Password) {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Invalid credentials"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	token, err := utils.CreateJWT(user.ID, user.Username, user.Email)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error generating token", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	sanitizedUser := sanitizedUserData(user)
	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Login successful", Data: map[string]any{"user": sanitizedUser}}
	c.JSON(http.StatusOK, res)
}

func (s *Server) DeleteUserById(c *gin.Context) {
	id, err := parseUintParam(c, "user_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid user ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := s.db.DeleteUser(id); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to delete user", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User deleted successfully"}
	c.JSON(http.StatusOK, res)
}

func (s *Server) UpdateUserById(c *gin.Context) {
	id, err := parseUintParam(c, "user_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid user ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}
	user.ID = id

	if err := s.db.UpdateUser(&user); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error updating user", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User updated successfully"}
	c.JSON(http.StatusOK, res)
}

// Blog

// GetAllBlogs handles retrieving all blogs
func (s *Server) GetAllBlogs(c *gin.Context) {
	blogs, err := s.db.GetBlogs()
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to fetch blogs", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blogs fetched successfully", Data: map[string]any{"blogs": blogs}}
	c.JSON(http.StatusOK, res)
}

// GetBlogByID handles retrieving a blog by ID
func (s *Server) GetBlogByID(c *gin.Context) {
	id, err := parseUintParam(c, "blog_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	blog, err := s.db.GetBlog(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error fetching blog", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if blog == nil {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Blog not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog fetched successfully", Data: map[string]any{"blog": blog}}
	c.JSON(http.StatusOK, res)
}

// CreateNewBlog handles creating a new blog
func (s *Server) CreateNewBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog data", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Set UserID from authenticated user (assuming middleware sets it)
	userID, exists := c.Get("user_id")
	if !exists {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "User not authenticated"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	blog.UserID = userID.(uint) // Ensure type casting is correct

	createdBlog, err := s.db.CreateBlog(&blog)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to create blog", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusCreated, Success: true, Message: "Blog created successfully", Data: map[string]any{"blog": createdBlog}}
	c.JSON(http.StatusCreated, res)
}

// DeleteBlogByID handles deleting a blog by ID
func (s *Server) DeleteBlogByID(c *gin.Context) {
	id, err := parseUintParam(c, "blog_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = s.db.DeleteBlog(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Blog not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to delete blog", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog deleted successfully"}
	c.JSON(http.StatusOK, res)
}

// UpdateBlog handles updating a blog
func (s *Server) UpdateBlog(c *gin.Context) {
	id, err := parseUintParam(c, "blog_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var input models.Blog
	if err := c.ShouldBindJSON(&input); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog data", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Fetch the existing blog
	existingBlog, err := s.db.GetBlog(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error fetching blog", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if existingBlog == nil {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Blog not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}

	// Only update allowed fields
	existingBlog.Title = input.Title
	existingBlog.Content = input.Content

	err = s.db.UpdateBlog(existingBlog)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to update blog", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog updated successfully", Data: map[string]any{"blog": existingBlog}}
	c.JSON(http.StatusOK, res)
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

// Comment handlers

// GetAllComments handles retrieving all comments for a blog
func (s *Server) GetAllComments(c *gin.Context) {
	blogID, err := parseUintParam(c, "blog_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	comments, err := s.db.GetComments(blogID)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to fetch comments", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comments fetched successfully", Data: map[string]any{"comments": comments}}
	c.JSON(http.StatusOK, res)
}

// GetCommentByID handles retrieving a comment by ID
func (s *Server) GetCommentByID(c *gin.Context) {
	id, err := parseUintParam(c, "comment_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	comment, err := s.db.GetComment(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error fetching comment", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if comment == nil {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Comment not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment fetched successfully", Data: map[string]any{"comment": comment}}
	c.JSON(http.StatusOK, res)
}

// CreateNewComment handles creating a new comment
func (s *Server) CreateNewComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment data", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := s.db.CreateComment(&comment); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to create comment", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusCreated, Success: true, Message: "Comment created successfully", Data: map[string]any{"comment": comment}}
	c.JSON(http.StatusCreated, res)
}

// UpdateComment handles updating a comment
func (s *Server) UpdateComment(c *gin.Context) {
	id, err := parseUintParam(c, "comment_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment data", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	comment.ID = id

	err = s.db.UpdateComment(&comment)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Comment not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to update comment", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment updated successfully", Data: map[string]any{"comment": comment}}
	c.JSON(http.StatusOK, res)
}

// DeleteCommentByID handles deleting a comment by ID
func (s *Server) DeleteCommentByID(c *gin.Context) {
	id, err := parseUintParam(c, "comment_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = s.db.DeleteComment(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Comment not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to delete comment", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment deleted successfully"}
	c.JSON(http.StatusOK, res)
}

// like handlers

// Like a blog
func (s *Server) LikeBlog(c *gin.Context) {
	var like models.Like
	if err := c.ShouldBindJSON(&like); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid like data", Error: err.Error()})
		return
	}

	if err := s.db.LikeBlog(&like); err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to like blog", Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, types.Response{StatusCode: http.StatusCreated, Success: true, Message: "Blog liked successfully"})
}

// Unlike a blog
func (s *Server) UnlikeBlog(c *gin.Context) {
	blogID, _ := parseUintParam(c, "blog_id")
	userID, _ := parseUintParam(c, "user_id")

	if err := s.db.UnlikeBlog(userID, blogID); err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to unlike blog", Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog unliked successfully"})
}
