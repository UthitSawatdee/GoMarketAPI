package usecase

import(
	port "github.com/Fal2o/E-Commerce_API/internal/port"
	domain "github.com/Fal2o/E-Commerce_API/internal/domain"
	"fmt"
	hash "github.com/Fal2o/E-Commerce_API/pkg/hash"

)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	CreateUser(user *domain.User) error
	// UpdateUser(user *domain.User) error
	// DeleteUser(id uint) error
}

type UserService struct {
    repo port.UserRepository
	hash hash.PasswordService

}

func NewUserService(repo port.UserRepository,hash hash.PasswordService) UserUseCase {
	return &UserService{
		repo: repo,
		hash: hash,
	}
}

func (s *UserService) CreateUser(user *domain.User) (error) {
	// 1. Check if email already exists
	existingUser, err := s.repo.GetByEmail(user.Email)
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
