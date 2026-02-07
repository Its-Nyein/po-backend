package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"po-backend/models"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBName        string
	DBUsername    string
	DBPassword    string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	JWTSecret     string
	DB            *gorm.DB
	Redis         *redis.Client
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	redisDB, _ := strconv.Atoi(getEnvOrDefault("REDIS_DB", "0"))

	return &Config{
		ServerPort:    getEnvOrDefault("SERVER_PORT", "8000"),
		DBHost:        getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:        getEnvOrDefault("DB_PORT", "5432"),
		DBName:        getEnvOrDefault("DB_NAME", "po"),
		DBUsername:    getEnvOrDefault("DB_USERNAME", "postgres"),
		DBPassword:    getEnvOrDefault("DB_PASSWORD", "postgres"),
		RedisHost:     getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     getEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "your-secret-key"),
	}
}

func (c *Config) ConnectDB() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, c.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	log.Println("Successfully connected to PostgreSQL")
	c.DB = db
	return nil
}

func (c *Config) ConnectRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.RedisHost + ":" + c.RedisPort,
		Password: c.RedisPassword,
		DB:       c.RedisDB,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Successfully connected to Redis")
	c.Redis = rdb
	return nil
}

func (c *Config) InitializeDB() error {
	if err := c.DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
		&models.PostLike{},
		&models.CommentLike{},
		&models.Follow{},
		&models.Notification{},
	); err != nil {
		return err
	}
	log.Println("Database migrations completed")
	return nil
}

var Envs = LoadConfig()
