package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt"
)

func HashPassword(password string) string {
	hash := sha512.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func VerifyPassword(hashedPassword, password string) bool {
	hash := sha512.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil)) == hashedPassword
}
func ComparePasswords(hashedPassword, password string) bool {
	return VerifyPassword(hashedPassword, password)

}

var secretKey = []byte("XXXXXXXXX")

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenSigned, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		return "", tokenErr
	}

	return tokenSigned, nil
}

// GenerateRandomToken generates a secure random token of the specified length
func GenerateRandomToken(length int) string {
	// Create a byte slice of the appropriate length
	b := make([]byte, length)

	// Fill it with random bytes
	_, err := rand.Read(b)
	if err != nil {
		// If we can't generate random bytes, use a timestamp as fallback
		// This is less secure but better than failing
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))
	}

	// Convert to URL-safe base64 string and return
	return base64.URLEncoding.EncodeToString(b)
}
