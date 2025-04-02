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
			protectedUser.POST("/follow/:user_id", s.FollowUser)
			protectedUser.DELETE("/unfollow/:user_id", s.UnfollowUser)
			protectedUser.POST("/logout", s.LogoutUser)
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
			blog.POST("/:blog_id/view", s.UpdateViewHandler)

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
			comment.DELETE("/", s.DeleteCommentByID)
			comment.GET("/:comment_id", s.GetCommentByID)
			comment.PUT("/:comment_id", s.UpdateComment)
		}
		// Admin Routes
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware()) // Ensure only admins can access
		{
			admin.GET("/dashboard", s.GetAdminDashboard) // Admin dashboard route
			admin.GET("/users", s.AdminGetUsers)         // Admin route to get all users
			admin.GET("/user/:id", s.AdminGetUser)       // Admin route to get a single user by ID
			admin.DELETE("/user/:id", s.AdminDeleteUser) // Admin route to delete a user
			admin.PUT("/user", s.AdminUpdateUser)        // Admin route to update a user

			admin.GET("/blogs", s.AdminGetBlogs)         // Admin route to get all blogs
			admin.GET("/blog/:id", s.AdminGetBlog)       // Admin route to get a single blog by ID
			admin.DELETE("/blog/:id", s.AdminDeleteBlog) // Admin route to delete a blog
			admin.PUT("/blog", s.AdminUpdateBlog)        // Admin route to update a blog

			admin.GET("/comments", s.AdminGetComments)         // Admin route to get all comments
			admin.DELETE("/comment/:id", s.AdminDeleteComment) // Admin route to delete a comment
			admin.PUT("/comment", s.AdminUpdateComment)        // Admin route to update a comment
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
