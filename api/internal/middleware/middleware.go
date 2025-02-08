package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateTeamToken(teamID int) (string, error) {
	return generateToken(teamID, "team")
}

func GenerateAdminToken(adminID int) (string, error) {
	return generateToken(adminID, "admin")
}

func generateToken(teamID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   teamID,
		"role": role, // Include role in the token
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	return token.SignedString([]byte("secret-key")) // Ensure consistent secret key
}


type InvalidTokenError struct {
	error
}

func (e *InvalidTokenError) Error() string {
	return "invalid token"
}

func ValidateToken(token string) (int, error) {
	Tokenvalue, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte("secret-key"), nil
	})
	if err != nil {
		return 0, &InvalidTokenError{err}
	}

	if Tokenvalue == nil || !Tokenvalue.Valid {
		return 0, &InvalidTokenError{fmt.Errorf("invalid token")}
	}

	var parsedID interface{}
	if claims, ok := Tokenvalue.Claims.(jwt.MapClaims); ok && Tokenvalue.Valid {
		parsedID = claims["id"]
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return 0, &InvalidTokenError{fmt.Errorf("token expired")}
		}
	}
	value, ok := parsedID.(float64)
	if !ok {
		return 0, &InvalidTokenError{fmt.Errorf("expected an int value, but got %T", parsedID)}
	}

	id := int(value)
	return id, nil
}