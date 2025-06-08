package utils

import (
	"crypto/sha512"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"sub": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenSigned, tokenErr := token.SignedString(secretKey)
	if tokenErr != nil {
		return "", tokenErr
	}

	return tokenSigned, nil
}