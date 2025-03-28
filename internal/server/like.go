package server

import (
	"net/http"
	"obs/internal/models"
	"obs/internal/types"
	"strings"

	"github.com/gin-gonic/gin"
)

// like handlers
// Like a blog
func (s *Server) LikeBlog(c *gin.Context) {
	// Extract user_id from JWT token (middleware should set this)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.Response{
			StatusCode: http.StatusUnauthorized,
			Success:    false,
			Message:    "Unauthorized: User ID not found",
		})
		return
	}

	// Parse blog_id from request body
	var request struct {
		BlogID uint `json:"blog_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid request body",
			Error:      err.Error(),
		})
		return
	}

	// Check if blog exists
	blog, err := s.db.GetBlog(request.BlogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Failed to fetch blog",
			Error:      err.Error(),
		})
		return
	}
	if blog == nil {
		c.JSON(http.StatusNotFound, types.Response{
			StatusCode: http.StatusNotFound,
			Success:    false,
			Message:    "Blog not found",
		})
		return
	}

	// Create like entry
	likeEntry := models.Like{
		UserID: userID.(uint),
		BlogID: request.BlogID,
	}

	// Attempt to insert like
	if err := s.db.LikeBlog(&likeEntry); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			c.JSON(http.StatusConflict, types.Response{
				StatusCode: http.StatusConflict,
				Success:    false,
				Message:    "You have already liked this blog",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, types.Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Failed to like blog",
			Error:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Blog liked successfully",
	})
}

// Unlike a blog
func (s *Server) UnlikeBlog(c *gin.Context) {
	// Extract user_id from JWT token (middleware should set this)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.Response{
			StatusCode: http.StatusUnauthorized,
			Success:    false,
			Message:    "Unauthorized: User ID not found",
		})
		return
	}

	// Parse blog_id from request body
	var request struct {
		BlogID uint `json:"blog_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid request body",
			Error:      err.Error(),
		})
		return
	}

	// Check if the like exists
	like, err := s.db.GetLikeByUserAndBlog(userID.(uint), request.BlogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Failed to fetch like",
			Error:      err.Error(),
		})
		return
	}
	if like == nil {
		c.JSON(http.StatusNotFound, types.Response{
			StatusCode: http.StatusNotFound,
			Success:    false,
			Message:    "You have not liked this blog",
		})
		return
	}

	// Call unlike function
	if err := s.db.UnlikeBlog(userID.(uint), request.BlogID); err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Failed to unlike blog",
			Error:      err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Blog unliked successfully",
	})
}
