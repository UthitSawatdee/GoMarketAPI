package usecase

import (
	"fmt"
	"log"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	hash "github.com/Fal2o/E-Commerce_API/pkg/hash"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	CreateUser(user *domain.User) error
	// UpdateUser(user *domain.User) error
	// DeleteUser(id uint) error
	LoginUser(user *domain.User) (error, string)
	GetUserByID(id uint) (*domain.User, error)
	UpdateUser(user *domain.User, password, newPassword string) error
}

type UserService struct {
	repo port.UserRepository
	hash hash.PasswordService
}

func NewUserService(repo port.UserRepository, hash hash.PasswordService) UserUseCase {
	return &UserService{
		repo: repo,
		hash: hash,
	}
}

func (s *UserService) CreateUser(user *domain.User) error {
	// 1. Check if email already exists
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return fmt.Errorf("email already registered")
	}

	// 2. Hash password
	hashedPassword, err := s.hash.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// 3. Create user
	return s.repo.Create(user)
}

func (s *UserService) LoginUser(user *domain.User) (error, string) {
	// 1. Retrieve user by email
	existingUser, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return err, ""
	}
	if existingUser == nil {
		return fmt.Errorf("invalid email or password"), ""
	}

	// 2. Verify password
	if err := s.hash.Verify(user.Password,existingUser.Password); !err {
		log.Println("Password verification failed:", err)
		return fmt.Errorf("invalid email or password"), ""
	}

	// 3. Generate JWT token , keep in local variable to avoid conflict with imported package
	_token := jwt.New(jwt.SigningMethodHS256)
	claims := _token.Claims.(jwt.MapClaims)
	claims["user_id"] = existingUser.ID
	claims["email"] = existingUser.Email
	claims["username"] = existingUser.Username
	claims["role"] = existingUser.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iat"] = time.Now().Unix()

	// 4. Sign token
	token, err := _token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err, ""
	}

	return nil, token
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUser(user *domain.User,password,newPassword string) error {
	if password == "" && newPassword == "" {
		return s.repo.Update(user)
	}
	fmt.Println("New Password to be set:", newPassword)
	existingUser, err := s.repo.GetUserByEmail(user.Email)

	// 1. Verify current password
	if !s.hash.Verify(password, existingUser.Password) {
		return fmt.Errorf("current password is incorrect")
	}
	fmt.Println("Current Password verified successfully")
	hashedPassword, err := s.hash.Hash(newPassword)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.Update(user)
}