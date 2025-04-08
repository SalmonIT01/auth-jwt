// cmd/main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"auth-jwt/config"
	handlers "auth-jwt/internal/handlers"
	"auth-jwt/internal/middleware"
	"auth-jwt/internal/repository"
	"auth-jwt/internal/services"
)

func main() {
	// โหลดการตั้งค่า
	cfg := config.LoadConfig()

	// เชื่อมต่อ MongoDB
	client, err := connectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	db := client.Database(cfg.MongoDB)

	// สร้างส่วนต่างๆ ของแอปพลิเคชัน
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	// ตั้งค่า Gin router
	router := gin.Default()

	// กำหนด routes
	router.POST("/api/register", userHandler.Register)
	router.POST("/api/login", userHandler.Login)

	// Route ที่ต้องการการยืนยันตัวตน
	authRoutes := router.Group("/api")
	authRoutes.Use(authMiddleware.AuthRequired())
	{
		authRoutes.GET("/profile", userHandler.GetProfile)
		authRoutes.PUT("/profile", userHandler.UpdateUser)    // เพิ่ม route สำหรับอัปเดตข้อมูล
		authRoutes.DELETE("/profile", userHandler.DeleteUser) // เพิ่ม route สำหรับลบบัญชี
	}

	// เริ่มเซิร์ฟเวอร์
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// ตรวจสอบการเชื่อมต่อ
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// go.mod (ตัวอย่าง)
// module auth-api
// go 1.19
// require (
//     github.com/gin-gonic/gin v1.9.1
//     github.com/golang-jwt/jwt/v5 v5.0.0
//     go.mongodb.org/mongo-driver v1.12.1
//     golang.org/x/crypto v0.14.0
// )
