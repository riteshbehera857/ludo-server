package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var config *Config

type Config struct {
	Port               string
	MongoURI           string
	Database           string
	BasePlatformAPIUrl string
}

func GetConfig() Config {
	if config == nil {
		LoadConfig()
	}
	return *config
}

func LoadConfig() {
	log.Println("Loading config")
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file : %s", err)
	}

	config = &Config{
		Port:               getEnv("PORT", "8080"),
		MongoURI:           getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:           getEnv("DATABASE", "gameserver"),
		BasePlatformAPIUrl: getEnv("BASE_PLATFORM_API_URL", "http://localhost:4000"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
