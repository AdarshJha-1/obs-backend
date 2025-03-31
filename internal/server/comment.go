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

// GetAllComments retrieves comments for a blog
func (s *Server) GetAllComments(c *gin.Context) {
	blogID, err := utils.ParseUintParam(c, "blog_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID"})
		return
	}

	comments, err := s.db.GetComments(blogID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error fetching comments"})
		return
	}

	c.JSON(http.StatusOK, types.Response{StatusCode: http.StatusOK, Success: true, Data: gin.H{"comments": comments}})
}

// GetCommentByID retrieves a single comment
func (s *Server) GetCommentByID(c *gin.Context) {
	commentID, err := utils.ParseUintParam(c, "comment_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment ID"})
		return
	}

	comment, err := s.db.GetComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error fetching comment"})
		return
	}
	if comment == nil {
		c.JSON(http.StatusNotFound, types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, types.Response{StatusCode: http.StatusOK, Success: true, Data: gin.H{"comment": comment}})
}

// CreateNewComment inserts a new comment
func (s *Server) CreateNewComment(c *gin.Context) {
	var input struct {
		BlogID  uint   `json:"blog_id" binding:"required"`
		Content string `json:"content" binding:"required,min=3"`
	}
	fmt.Printf("%+v\n", input)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input data"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized"})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized"})
		return
	}
	comment := models.Comment{
		BlogID:  input.BlogID,
		UserID:  userID.(uint),
		Author:  username.(string),
		Content: input.Content,
	}

	if err := s.db.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, types.Response{StatusCode: http.StatusCreated, Success: true, Data: gin.H{"comment": comment}})
}

// UpdateComment modifies a comment if the user is the owner
func (s *Server) UpdateComment(c *gin.Context) {
	commentID, err := utils.ParseUintParam(c, "comment_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid comment ID"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required,min=3"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid content"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized"})
		return
	}

	err = s.db.UpdateComment(commentID, userID.(uint), input.Content)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Comment not found or not owned by user"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment updated successfully"})
}

// DeleteCommentByID deletes a comment if the user is the owner
func (s *Server) DeleteCommentByID(c *gin.Context) {
	var input struct {
		CommentId uint `json:"comment_id" binding:"required"`
		UserId    uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input data"})
		return
	}
	err := s.db.DeleteComment(input.CommentId, input.UserId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, types.Response{StatusCode: http.StatusNotFound, Success: false, Message: "Comment not found or not owned by user"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, types.Response{StatusCode: http.StatusOK, Success: true, Message: "Comment deleted successfully"})
}
