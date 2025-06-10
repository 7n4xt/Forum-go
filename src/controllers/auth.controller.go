package controllers

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"forum-go/models"
	"forum-go/repositories"
	"forum-go/services"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderRegistrationForm(w, r, nil)
	case "POST":
		handleRegistration(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderLoginForm(w, r, nil)
	case "POST":
		handleLogin(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func renderRegistrationForm(w http.ResponseWriter, r *http.Request, errorMsg *string) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Error *string
	}{
		Error: errorMsg,
	}

	tmpl.Execute(w, data)
}

func renderLoginForm(w http.ResponseWriter, r *http.Request, errorMsg *string) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Error *string
	}{
		Error: errorMsg,
	}

	tmpl.Execute(w, data)
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		errorMsg := "Error parsing form data"
		renderRegistrationForm(w, r, &errorMsg)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Validate input
	if validationError := validateRegistrationInput(username, email, password, confirmPassword); validationError != nil {
		renderRegistrationForm(w, r, validationError)
		return
	}

	// Check if user already exists
	existingUser, err := repositories.ReadUserByUsername(username)
	if err == nil && existingUser != nil {
		errorMsg := "Username already exists"
		renderRegistrationForm(w, r, &errorMsg)
		return
	}

	existingUser, err = repositories.ReadUserByEmail(email)
	if err == nil && existingUser != nil {
		errorMsg := "Email already exists"
		renderRegistrationForm(w, r, &errorMsg)
		return
	}

	// Hash password
	passwordHash := hashPassword(password)

	// Create user
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "user",
		CreatedAt:    time.Now(),
		IsBanned:     false,
	}

	userID, err := repositories.CreateUser(*user)
	if err != nil {
		errorMsg := "Error creating user account"
		renderRegistrationForm(w, r, &errorMsg)
		return
	}

	// Generate session token
	token, err := services.GenerateSessionToken(userID)
	if err != nil {
		errorMsg := "Error creating session"
		renderRegistrationForm(w, r, &errorMsg)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		errorMsg := "Error parsing form data"
		renderLoginForm(w, r, &errorMsg)
		return
	}

	usernameOrEmail := strings.TrimSpace(r.FormValue("username_email"))
	password := r.FormValue("password")

	// Validate input
	if usernameOrEmail == "" || password == "" {
		errorMsg := "Username/Email and password are required"
		renderLoginForm(w, r, &errorMsg)
		return
	}

	// Get user by username or email
	var user *models.User
	if isEmail(usernameOrEmail) {
		user, err = repositories.ReadUserByEmail(usernameOrEmail)
	} else {
		user, err = repositories.ReadUserByUsername(usernameOrEmail)
	}

	if err != nil || user == nil {
		errorMsg := "Invalid username/email or password"
		renderLoginForm(w, r, &errorMsg)
		return
	}

	// Check if user is banned
	if user.IsBanned {
		errorMsg := "Your account has been banned"
		renderLoginForm(w, r, &errorMsg)
		return
	}

	// Verify password
	if !verifyPassword(password, user.PasswordHash) {
		errorMsg := "Invalid username/email or password"
		renderLoginForm(w, r, &errorMsg)
		return
	}

	// Update last login
	repositories.UpdateLastLogin(user.UserId)

	// Generate session token
	token, err := services.GenerateSessionToken(user.UserId)
	if err != nil {
		errorMsg := "Error creating session"
		renderLoginForm(w, r, &errorMsg)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   3600 * 24 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func validateRegistrationInput(username, email, password, confirmPassword string) *string {
	if username == "" || email == "" || password == "" {
		msg := "All fields are required"
		return &msg
	}

	if len(username) < 3 || len(username) > 50 {
		msg := "Username must be between 3 and 50 characters"
		return &msg
	}

	if !isValidEmail(email) {
		msg := "Invalid email format"
		return &msg
	}

	if len(password) < 6 {
		msg := "Password must be at least 6 characters long"
		return &msg
	}

	if password != confirmPassword {
		msg := "Passwords do not match"
		return &msg
	}

	// Check for valid username characters (alphanumeric and underscore only)
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		msg := "Username can only contain letters, numbers, and underscores"
		return &msg
	}

	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

func isEmail(input string) bool {
	return strings.Contains(input, "@")
}

func hashPassword(password string) string {
	hasher := sha512.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

// Helper function to render JSON responses (for future AJAX endpoints)
func renderJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
