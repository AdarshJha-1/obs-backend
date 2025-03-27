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
			user.POST("/login", s.LoginUser)
			user.DELETE("/u/:user_id", s.DeleteUserById)
			user.PUT("/u/:user_id", s.UpdateUserById)
		}
		blog := api.Group("/blog")
		{
			blog.GET("/all", s.GetAllBlogs)
			blog.POST("/", s.CreateNewBlog)
			blog.GET("/b/:blog_id", s.GetBlogByID)
			blog.DELETE("/b/:blog_id", s.DeleteBlogByID)
			blog.PUT("/b/:blog_id", s.UpdateBlog)
			// Nested Comments under a Blog
			comments := blog.Group("/:blog_id/comments")
			{
				comments.GET("/", s.GetAllComments)    // Get all comments for a blog
				comments.POST("/", s.CreateNewComment) // Add a new comment to a blog
			}
		}
		comment := api.Group("/comment")
		{
			comment.GET("/:comment_id", s.GetCommentByID)       // Get a single comment by ID
			comment.PUT("/:comment_id", s.UpdateComment)        // Update a comment by ID
			comment.DELETE("/:comment_id", s.DeleteCommentByID) // Delete a comment by ID
		}
		like := api.Group("/like")
		{
			like.GET("/:like_id", s.LikeBlog)      // Get a like by ID
			like.DELETE("/:like_id", s.UnlikeBlog) // Remove a like by ID
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
