package models

import "golang.org/x/crypto/bcrypt"

type Team struct {
	ID           int   `gorm:"primaryKey"`
	Name         string `gorm:"unique"`
	PasswordHash string
	Score        int
	IsBlocked    bool
	IsAdmin      bool `gorm:"default:false"` // New field
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
