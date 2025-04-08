// internal/models/user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"` // ไม่ส่งรหัสผ่านกลับใน JSON
	Fullname  string             `bson:"fullname" json:"fullname"`
	Tel       string             `bson:"tel" json:"tel"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserResponse struct {
	ID       primitive.ObjectID `json:"id,omitempty"`
	Username string             `json:"username"`
	Fullname string             `json:"fullname"`
	Tel      string             `json:"tel"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Tel      string `json:"tel"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresIn int64        `json:"expires_in"`
}

// สำหรับอัปเดตข้อมูลผู้ใช้
type UpdateUserRequest struct {
	Fullname string `json:"fullname,omitempty"`
	Tel      string `json:"tel,omitempty"`
	Password string `json:"password,omitempty"` // สำหรับอัปเดตรหัสผ่าน (ถ้าต้องการ)
}
