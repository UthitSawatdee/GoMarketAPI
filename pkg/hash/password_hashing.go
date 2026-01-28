package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Interface - กำหนดว่า service นี้ทำอะไรได้บ้าง
type PasswordService interface {
    Hash(password string) (string, error)
    Verify(password, hashedPassword string) bool
}

// Implementation - วิธีทำจริง (ใช้ bcrypt)
type BcryptPasswordService struct{}

func NewPasswordService() PasswordService {
    return &BcryptPasswordService{}
}

func (s *BcryptPasswordService) Hash(password string) (string, error) {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashed), err
}

func (s *BcryptPasswordService) Verify(password, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    fmt.Println("Password verification error:", err)
    return err == nil
}