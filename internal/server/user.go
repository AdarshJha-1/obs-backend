package server

import (
	"net/http"
	"obs/internal/models"
	"obs/internal/types"
	"obs/internal/utils"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized access"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	id := userID.(uint)
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
		Data:       map[string]any{"user": utils.SanitizedUserData(user)},
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) GetUserById(c *gin.Context) {
	id, err := utils.ParseUintParam(c, "user_id")
	if err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid blog ID", Error: err.Error()}
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
		Data:       map[string]any{"user": utils.SanitizedUserData(user)},
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

	sanitizedUsers := make([]utils.SanitizedUser, len(users))
	for i, user := range users {
		sanitizedUsers[i] = utils.SanitizedUserData(&user)
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "Users fetched successfully", Data: map[string]any{"users": sanitizedUsers}}
	c.JSON(http.StatusOK, res)
}

func (s *Server) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized access"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		res := types.Response{StatusCode: http.StatusBadRequest, Success: false, Message: "Invalid input", Error: err.Error()}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	user.ID = userID.(uint)

	if err := s.db.UpdateUser(&user); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Error updating user", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User updated successfully"}
	c.JSON(http.StatusOK, res)
}

func (s *Server) DeleteCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		res := types.Response{StatusCode: http.StatusUnauthorized, Success: false, Message: "Unauthorized access"}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if err := s.db.DeleteUser(userID.(uint)); err != nil {
		res := types.Response{StatusCode: http.StatusInternalServerError, Success: false, Message: "Failed to delete user", Error: err.Error()}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res := types.Response{StatusCode: http.StatusOK, Success: true, Message: "User deleted successfully"}
	c.JSON(http.StatusOK, res)
}
