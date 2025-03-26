package server

import (
	"fmt"
	"net/http"
	"obs/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetUser(c *gin.Context) {
	strId := c.Param("user_id")
	id, err := strconv.ParseUint(strId, 10, 64)
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
	var User models.User

	if err := c.ShouldBindJSON(&User); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	fmt.Println("user: ", User.Password)
	if err := s.db.CreateUser(&User); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating user"})
	}

}

func (s *Server) DeleteUserById(c *gin.Context) {
	strId := c.Param("user_id")
	id, err := strconv.ParseUint(strId, 10, 64)

	if err = s.db.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "user deleted successfully"})
}
func (s *Server) UpdateUserById(c *gin.Context) {
	strId := c.Param("user_id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var user models.User

	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user.ID = id
	if err = s.db.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "user updated successfully"})
}

// Blog
func (s *Server) GetAllBlogs(c *gin.Context) {
	blogs, err := s.db.GetBlogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

func (s *Server) CreateNewBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	if err := s.db.CreateBlog(blog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"suucess": "blog created successfully"})
}

func (s *Server) GetBlogById(c *gin.Context) {
	strId := c.Param("blog_id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	blog, err := s.db.GetBlog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusFound, gin.H{"blog": blog})
}

func (s *Server) DeleteBlogById(c *gin.Context) {
	strId := c.Param("blog_id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	err = s.db.DeleteBlog(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusFound, gin.H{"success": "blog deleted successfully"})
}

func (s *Server) UpdateBlog(c *gin.Context) {
	strId := c.Param("blog_id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	blog.ID = id
	err = s.db.UpdateBlog(blog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"success": "blog updated successfully"})
}
