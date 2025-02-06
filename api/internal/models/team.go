package models

import "golang.org/x/crypto/bcrypt"

type Team struct {
	ID           uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	IsBlocked    bool   `gorm:"default:false"`
}

// HashPassword hashes the given password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a plaintext password with a hashed password
func CheckPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}


