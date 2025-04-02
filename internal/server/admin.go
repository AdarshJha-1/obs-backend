package server

import (
	"net/http"
	"obs/internal/models"
	"obs/internal/types"
	"obs/internal/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetUsers retrieves all users (admin access only)
func (s *Server) AdminGetUsers(c *gin.Context) {
	users, err := s.db.AdminGetUsers()
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Users retrieved successfully", Data: map[string]any{"users": users}}
	c.JSON(http.StatusOK, res)
}

// AdminGetUser retrieves a single user by ID (admin access only)
func (s *Server) AdminGetUser(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid user ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := s.db.AdminGetUser(id)
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

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User retrieved successfully", Data: map[string]any{"user": user}}
	c.JSON(http.StatusOK, res)
}

// AdminDeleteUser deletes a user by ID (admin access only)
func (s *Server) AdminDeleteUser(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid user ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = s.db.AdminDeleteUser(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User deleted successfully"}
	c.JSON(http.StatusOK, res)
}

// AdminUpdateUser updates an existing user's details (admin access only)
func (s *Server) AdminUpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err := s.db.AdminUpdateUser(&user)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User updated successfully", Data: map[string]any{"user": user}}
	c.JSON(http.StatusOK, res)
}

// AdminGetBlogs retrieves all blogs (admin access only)
func (s *Server) AdminGetBlogs(c *gin.Context) {
	blogs, err := s.db.AdminGetBlogs()
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blogs retrieved successfully", Data: map[string]any{"blogs": blogs}}
	c.JSON(http.StatusOK, res)
}

// AdminGetBlog retrieves a single blog by ID (admin access only)
func (s *Server) AdminGetBlog(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	blog, err := s.db.AdminGetBlog(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	if blog == nil {
		res := types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Blog not found"}
		c.JSON(http.StatusNotFound, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog retrieved successfully", Data: map[string]any{"blog": blog}}
	c.JSON(http.StatusOK, res)
}

// AdminDeleteBlog deletes a blog by ID (admin access only)
func (s *Server) AdminDeleteBlog(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = s.db.AdminDeleteBlog(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog deleted successfully"}
	c.JSON(http.StatusOK, res)
}

// AdminUpdateBlog updates a blog's information (admin access only)
func (s *Server) AdminUpdateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err := s.db.AdminUpdateBlog(&blog)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog updated successfully", Data: map[string]any{"blog": blog}}
	c.JSON(http.StatusOK, res)
}

// AdminDeleteComment deletes a comment by ID (admin access only)
func (s *Server) AdminDeleteComment(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment ID", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err = s.db.AdminDeleteComment(id)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment deleted successfully"}
	c.JSON(http.StatusOK, res)
}

// AdminUpdateComment updates a comment's content (admin access only)
func (s *Server) AdminUpdateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	err := s.db.AdminUpdateComment(comment.ID, comment.Content)
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Database error", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment updated successfully", Data: map[string]any{"comment": comment}}
	c.JSON(http.StatusOK, res)
}

// AdminGetComments retrieves all comments (admin access only)
func (s *Server) AdminGetComments(c *gin.Context) {
	// Call the database function to retrieve the comments
	comments, err := s.db.AdminGetComments()
	if err != nil {
		// Handle database error
		res := types.Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Database error",
			Error:      err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Return a successful response with the comments data
	res := types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Comments retrieved successfully",
		Data:       map[string]any{"comments": comments},
	}
	c.JSON(http.StatusOK, res)
}

// GetAdminDashboard retrieves the data for the admin dashboard (admin access only)
func (s *Server) GetAdminDashboard(c *gin.Context) {
	// Retrieve dashboard data from the database
	dashboardData, err := s.db.GetAdminDashboardData()
	if err != nil {
		// Handle database error
		res := types.Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Database error",
			Error:      err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	// Return a successful response with the dashboard data
	res := types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Dashboard data retrieved successfully",
		Data:       map[string]any{"dashboard": dashboardData},
	}
	c.JSON(http.StatusOK, res)
}
