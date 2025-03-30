package server

import (
	"net/http"
	"obs/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	// Public Routes
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	api := r.Group("/api")
	{
		// Public User Routes
		public := api.Group("/")
		{
			public.POST("/register", s.RegisterUser) // Public Route
			public.POST("/login", s.LoginUser)       // Public Route
		}

		// Protected User Routes
		protectedUser := api.Group("/user")
		protectedUser.Use(middleware.AuthMiddleware()) // Apply middleware separately
		{
			protectedUser.GET("/", s.GetCurrentUser)
			protectedUser.GET("/all", s.GetUsers)
			protectedUser.GET("/:user_id", s.GetUserById)
			protectedUser.DELETE("/", s.DeleteCurrentUser)
			protectedUser.PUT("/", s.UpdateCurrentUser)
		}

		// Protected Blog Routes
		blog := api.Group("/blog")
		blog.Use(middleware.AuthMiddleware()) // Apply middleware separately
		{
			blog.GET("/all", s.GetAllBlogs)
			blog.POST("/", s.CreateNewBlog)
			blog.GET("/b/:blog_id", s.GetBlogByID)
			blog.DELETE("/b/:blog_id", s.DeleteBlogByID)
			blog.PUT("/b/:blog_id", s.UpdateBlog)

			blog.POST("/like", s.LikeBlog)
			blog.DELETE("/unlike", s.UnlikeBlog)

			// Nested Comments under a Blog
			comments := blog.Group("/:blog_id/comments")
			{
				comments.GET("/", s.GetAllComments)
				comments.POST("/", s.CreateNewComment)
			}
		}

		// Protected Comment Routes
		comment := api.Group("/comment")
		comment.Use(middleware.AuthMiddleware()) // Apply middleware separately
		{
			comment.GET("/:comment_id", s.GetCommentByID)
			comment.PUT("/:comment_id", s.UpdateComment)
			comment.DELETE("/:comment_id", s.DeleteCommentByID)
		}

	}
	return r
}

// Public handlers
func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := map[string]string{"message": "Hello World"}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
