package utils

import (
	"errors"
	"forum-go/models"
	"forum-go/repositories"
	"net/http"
	"strconv"
)

// GetUserFromRequest retrieves the user from the request's cookie
func GetUserFromRequest(r *http.Request) (*models.User, error) {
	// Get token from cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	// Get session from token
	session, err := repositories.GetSessionByToken(cookie.Value)
	if err != nil {
		return nil, err
	}

	// If session is nil, return nil user
	if session == nil {
		return nil, errors.New("invalid session")
	}

	// Get user from session
	user, err := repositories.ReadUserById(session.UserId)
	if err != nil {
		return nil, err
	}

	// If user is nil, return nil user
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	_, err := GetUserFromRequest(r)
	return err == nil
}

// RequireAuthentication is a middleware that redirects to login if not authenticated
func RequireAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

// GetUserIdFromRequest gets the user ID from the request's cookie
func GetUserIdFromRequest(r *http.Request) (int, error) {
	user, err := GetUserFromRequest(r)
	if err != nil {
		return 0, err
	}
	return user.UserId, nil
}

// GetUserIdFromString converts a string to a user ID
func GetUserIdFromString(idStr string) (int, error) {
	if idStr == "" {
		return 0, errors.New("user ID is empty")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}

	return id, nil
}
