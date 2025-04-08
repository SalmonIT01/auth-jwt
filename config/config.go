// config/config.go
package config

import (
	"os"
	"time"
)

type Config struct {
	MongoURI     string
	MongoDB      string
	JWTSecret    string
	JWTExpiresIn time.Duration
	ServerPort   string
}

func LoadConfig() *Config {
	// สามารถใช้ environment variables หรือ configuration file
	// ในตัวอย่างนี้จะใช้ค่าตั้งต้นง่ายๆ
	return &Config{
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("MONGO_DB", "auth_db"),
		JWTSecret:    getEnv("JWT_SECRET", "your_secret_key"),
		JWTExpiresIn: 1 * time.Hour, // token หมดอายุใน 24 ชั่วโมง
		ServerPort:   getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
