// internal/services/user_service.go
package services

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"auth-jwt/config"
	"auth-jwt/internal/models"
	"auth-jwt/internal/repository"
	"auth-jwt/pkg/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

func NewUserService(userRepo *repository.UserRepository, config *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *UserService) Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	// ตรวจสอบข้อมูลที่จำเป็น
	if req.Username == "" || req.Password == "" || req.Fullname == "" || req.Tel == "" {
		return nil, errors.New("all fields are required")
	}

	// เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// สร้างผู้ใช้ใหม่
	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Fullname: req.Fullname,
		Tel:      req.Tel,
	}

	// บันทึกผู้ใช้ลงในฐานข้อมูล
	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// สร้าง token
	token, expiresAt, err := utils.GenerateToken(createdUser.ID, createdUser.Username, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, err
	}

	// สร้าง response
	response := &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:       createdUser.ID,
			Username: createdUser.Username,
			Fullname: createdUser.Fullname,
			Tel:      createdUser.Tel,
		},
		ExpiresIn: expiresAt.Unix(),
	}

	return response, nil
}

func (s *UserService) Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// ตรวจสอบข้อมูลที่จำเป็น
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// ค้นหาผู้ใช้จาก username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// ตรวจสอบรหัสผ่าน
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// สร้าง token
	token, expiresAt, err := utils.GenerateToken(user.ID, user.Username, s.config.JWTSecret, s.config.JWTExpiresIn)
	if err != nil {
		return nil, err
	}

	// สร้าง response
	response := &models.AuthResponse{
		Token: token,
		User: models.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Fullname: user.Fullname,
			Tel:      user.Tel,
		},
		ExpiresIn: expiresAt.Unix(),
	}

	return response, nil
}

func (s *UserService) GetProfile(userID string) (*models.UserResponse, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Fullname: user.Fullname,
		Tel:      user.Tel,
	}, nil
}

// อัปเดตข้อมูลผู้ใช้
func (s *UserService) UpdateUser(userID string, req models.UpdateUserRequest) (*models.UserResponse, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// สร้าง bson.M สำหรับเก็บข้อมูลที่จะอัปเดต
	updateData := bson.M{}

	// ตรวจสอบข้อมูลแต่ละรายการ ถ้ามีการส่งมาก็จะอัปเดต
	if req.Fullname != "" {
		updateData["fullname"] = req.Fullname
	}

	if req.Tel != "" {
		updateData["tel"] = req.Tel
	}

	// ถ้ามีการส่งรหัสผ่านมาด้วย ให้เข้ารหัสก่อนบันทึก
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updateData["password"] = string(hashedPassword)
	}

	// ถ้าไม่มีข้อมูลที่จะอัปเดต
	if len(updateData) == 0 {
		return nil, errors.New("no data to update")
	}

	// ทำการอัปเดตในฐานข้อมูล
	err = s.userRepo.Update(id, updateData)
	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลผู้ใช้ล่าสุดหลังอัปเดต
	updatedUser, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// สร้าง response
	return &models.UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Fullname: updatedUser.Fullname,
		Tel:      updatedUser.Tel,
	}, nil
}

// ลบบัญชีผู้ใช้
func (s *UserService) DeleteUser(userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	return s.userRepo.Delete(id)
}
