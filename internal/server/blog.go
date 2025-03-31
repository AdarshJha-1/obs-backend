package server

import (
	"errors"
	"fmt"
	"net/http"
	"obs/internal/models"
	"obs/internal/types"
	"obs/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
	id, err := utils.ParseUintParam(c, "blog_id")
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

	user := utils.SanitizedUserData(&blog.User)
	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Blog fetched successfully", Data: map[string]any{"blog": blog, "user": user}}
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

	author, exists := c.Get("username")
	if !exists {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "User not authenticated"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}
	blog.UserID = userID.(uint)
	blog.Author = author.(string)

	fmt.Printf("blog data: %+v\n", blog)
	fmt.Printf("Creating blog with Author: %s\n", blog.Author)
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
	id, err := utils.ParseUintParam(c, "blog_id")
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
	id, err := utils.ParseUintParam(c, "blog_id")
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

// UpdateViewHandler handles tracking unique blog views
func (s *Server) UpdateViewHandler(c *gin.Context) {
	// Extract user ID from middleware (must be authenticated)
	userID, exists := c.Get("user_id")
	if !exists {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "User not authenticated"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	// Extract blog ID from URL parameter
	blogID, err := utils.ParseUintParam(c, "blog_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID"}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Call the database function to track the view
	err = s.db.UpdateView(blogID, userID.(uint))
	if err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Could not update view", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "View updated successfully"}
	c.JSON(http.StatusOK, res)
}
