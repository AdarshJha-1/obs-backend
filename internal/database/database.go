package database

import (
	"context"
	"fmt"
	"log"
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
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
	schema   = os.Getenv("DB_SCHEMA")
)

type service struct {
	DB *gorm.DB
}

// ConnectDatabase initializes the GORM database connection
func New() Service {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable search_path=%s",
		host, username, password, database, port, schema)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Log SQL queries
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	service := &service{
		DB: db,
	}
	// Store DB instance globally

	log.Println("Connected to database successfully!")

	// Run AutoMigrations (This creates tables based on models)
	service.MigrateSchema()

	return service
}

// MigrateSchema runs auto-migrations for all models
func (s *service) MigrateSchema() {
	err := s.DB.AutoMigrate(&User{}, &Blog{}, &Comment{}, &Like{}, &Follow{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migrated successfully!")
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
	log.Println("Disconnected from database.")
	return sqlDB.Close()
}
