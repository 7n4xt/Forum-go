package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"forum-go/config"
	"time"
)

// GenerateSessionToken creates a new session token for a user
func GenerateSessionToken(userID int) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("Error generating token: %v", err)
	}

	token := hex.EncodeToString(tokenBytes)

	// Store token in database
	query := "INSERT INTO user_sessions (user_id, token, created_at) VALUES (?, ?, ?)"
	_, err = config.DbContext.Exec(query, userID, token, time.Now())
	if err != nil {
		return "", fmt.Errorf("Error storing session: %v", err)
	}

	return token, nil
}

// ValidateSessionToken checks if a session token is valid and returns the user ID
func ValidateSessionToken(token string) (int, error) {
	var userID int
	var createdAt time.Time

	query := "SELECT user_id, created_at FROM user_sessions WHERE token = ?"
	err := config.DbContext.QueryRow(query, token).Scan(&userID, &createdAt)
	if err != nil {
		return 0, fmt.Errorf("Invalid session token")
	}

	// Check if token is expired (7 days)
	if time.Since(createdAt) > 7*24*time.Hour {
		// Remove expired token
		DeleteSessionToken(token)
		return 0, fmt.Errorf("Session expired")
	}

	return userID, nil
}

// DeleteSessionToken removes a session token from the database
func DeleteSessionToken(token string) error {
	query := "DELETE FROM user_sessions WHERE token = ?"
	_, err := config.DbContext.Exec(query, token)
	if err != nil {
		return fmt.Errorf("Error deleting session: %v", err)
	}
	return nil
}

// CleanExpiredSessions removes old session tokens
func CleanExpiredSessions() error {
	query := "DELETE FROM user_sessions WHERE created_at < ?"
	expiredTime := time.Now().Add(-7 * 24 * time.Hour)
	_, err := config.DbContext.Exec(query, expiredTime)
	if err != nil {
		return fmt.Errorf("Error cleaning expired sessions: %v", err)
	}
	return nil
}