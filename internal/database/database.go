package database

import (
	"context"
	"fmt"
	"log"
	"obs/internal/models"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service interface {
	MigrateSchema()
	Health() map[string]string
	Close() error

	// User Methods
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUser(id uint) (*models.User, error)
	GetUsers() ([]models.User, error)
	DeleteUser(id uint) error
	UpdateUser(user *models.User) error
	FollowUser(followerID, followedID uint) error
	UnfollowUser(followerID, followedID uint) error
	IsFollowing(followerID, followedID uint) (bool, error)
	// Blog Methods
	GetBlogs() ([]models.Blog, error)
	GetBlog(id uint) (*models.Blog, error)
	CreateBlog(blog *models.Blog) (*models.Blog, error)
	UpdateBlog(blog *models.Blog) error
	DeleteBlog(id uint) error

	// Comment Methods
	GetComments(blogID uint) ([]models.Comment, error)
	GetComment(id uint) (*models.Comment, error)
	CreateComment(comment *models.Comment) error
	UpdateComment(id uint, userID uint, content string) error
	DeleteComment(id uint, userID uint) error
	UpdateView(blogId, userId uint) error

	// Like functions
	GetLikesForBlog(blogID uint) (int64, error)
	GetLikeByID(likeID uint) (*models.Like, error)
	GetLikeByUserAndBlog(userID, blogID uint) (*models.Like, error)
	LikeBlog(like *models.Like) error
	UnlikeBlog(userID, blogID uint) error
	DeleteLike(likeID uint) error

	// Admin functions
	// User-related methods
	AdminGetUsers() ([]models.User, error)
	AdminGetUser(id uint) (*models.User, error)
	AdminDeleteUser(id uint) error
	AdminUpdateUser(user *models.User) error

	// Blog-related methods
	AdminGetBlogs() ([]models.Blog, error)
	AdminGetBlog(id uint) (*models.Blog, error)
	AdminDeleteBlog(id uint) error
	AdminUpdateBlog(blog *models.Blog) error

	// Comment-related methods
	AdminDeleteComment(id uint) error
	AdminUpdateComment(id uint, content string) error
	AdminGetComments() ([]models.Comment, error)
	GetAdminDashboardData() (DashboardData, error)
}

var (
	database = getEnv("DB_DATABASE", "default_db")
	password = getEnv("DB_PASSWORD", "password")
	username = getEnv("DB_USERNAME", "postgres")
	port     = getEnv("DB_PORT", "5432")
	host     = getEnv("DB_HOST", "localhost")
	schema   = getEnv("DB_SCHEMA", "public")
)

type service struct {
	DB *gorm.DB
}

// New initializes the database connection
func New() Service {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		host, username, password, database, port, schema,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // Show only warnings
	})
	if err != nil {
		log.Fatalf("[DATABASE] ❌ Failed to connect: %v", err)
	}

	// Set database connection pool settings
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	}

	service := &service{DB: db}
	log.Println("[DATABASE] ✅ Connected successfully!")

	// Run AutoMigrations
	service.MigrateSchema()

	return service
}

// MigrateSchema runs auto-migrations for all models
func (s *service) MigrateSchema() {
	err := s.DB.AutoMigrate(&models.User{}, &models.Blog{}, &models.Comment{}, &models.Like{}, &models.Follow{}, &models.View{})
	if err != nil {
		log.Fatalf("[DATABASE] ❌ Migration failed: %v", err)
	}
	log.Println("[DATABASE] ✅ Migration successful!")
}

// Health checks the database connection status
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.DB.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("Failed to get DB instance: %v", err)
		return stats
	}

	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("DB down: %v", err)
		return stats
	}

	// Collect DB stats
	dbStats := sqlDB.Stats()
	stats["status"] = "up"
	stats["message"] = "Database is healthy"
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()

	return stats
}

// Close closes the database connection
func (s *service) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}
	log.Println("[DATABASE] ❌ Disconnected.")
	return sqlDB.Close()
}

// getEnv fetches environment variables with a default fallback
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Printf("[WARNING] ⚠️ Missing env: %s, using default: %s", key, fallback)
		return fallback
	}
	return value
}
