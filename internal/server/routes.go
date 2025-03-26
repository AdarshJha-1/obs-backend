package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.GET("/u/:user_id", s.GetUser)
			user.GET("/all", s.GetUsers)
			user.POST("/register", s.RegisterUser)
			user.DELETE("/u/:user_id", s.DeleteUserById)
			user.PUT("/u/:user_id", s.UpdateUserById)
		}
		blog := api.Group("/blog")
		{
			blog.GET("/", s.GetAllBlogs)
			blog.POST("/", s.CreateNewBlog)
			blog.GET("/b/:blog_id", s.GetBlogById)
			blog.DELETE("/b/:blog_id", s.DeleteBlogById)
			blog.PUT("/b/:blog_id", s.UpdateBlog)
		}
	}
	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
