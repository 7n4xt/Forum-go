package repositories

import (
	"database/sql"
	"fmt"
	"forum-go/config"
	"time"
)

// Session represents a user session in the database
type Session struct {
	SessionId int
	UserId    int
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// CreateSession inserts a new session into the database
func CreateSession(userId int, token string, expiresAt time.Time) (*Session, error) {
	db := config.DbContext

	// Insert session
	result, err := db.Exec(
		"INSERT INTO user_sessions (user_id, token, expires_at) VALUES (?, ?, ?)",
		userId, token, expiresAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating session: %w", err)
	}

	// Get the inserted session ID
	sessionId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting session ID: %w", err)
	}

	// Return session
	return &Session{
		SessionId: int(sessionId),
		UserId:    userId,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}, nil
}

// GetSessionByToken retrieves a session by its token
func GetSessionByToken(token string) (*Session, error) {
	db := config.DbContext

	// Query for session
	var session Session
	err := db.QueryRow(
		"SELECT session_id, user_id, token, created_at, expires_at FROM user_sessions WHERE token = ? AND expires_at > ?",
		token, time.Now(),
	).Scan(
		&session.SessionId,
		&session.UserId,
		&session.Token,
		&session.CreatedAt,
		&session.ExpiresAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Session not found or expired
		}
		return nil, fmt.Errorf("error getting session: %w", err)
	}

	return &session, nil
}

// DeleteSession removes a session from the database
func DeleteSession(token string) error {
	db := config.DbContext

	// Delete session
	_, err := db.Exec("DELETE FROM user_sessions WHERE token = ?", token)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions from the database
func DeleteExpiredSessions() (int64, error) {
	db := config.DbContext

	// Delete expired sessions
	result, err := db.Exec("DELETE FROM user_sessions WHERE expires_at <= ?", time.Now())
	if err != nil {
		return 0, fmt.Errorf("error deleting expired sessions: %w", err)
	}

	// Get the number of deleted sessions
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting deleted session count: %w", err)
	}

	return count, nil
}

// DeleteSessionsByUserId removes all sessions for a user
func DeleteSessionsByUserId(userId int) (int64, error) {
	db := config.DbContext

	// Delete sessions
	result, err := db.Exec("DELETE FROM user_sessions WHERE user_id = ?", userId)
	if err != nil {
		return 0, fmt.Errorf("error deleting user sessions: %w", err)
	}

	// Get the number of deleted sessions
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error getting deleted session count: %w", err)
	}

	return count, nil
}
