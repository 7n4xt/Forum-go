package services

import (
	"fmt"
	"forum-go/models"
	"forum-go/repositories"
	"forum-go/utils"
	"time"
)

func SignupService(username, password, email string) (bool, error) {
	// Here you would typically interact with a database to save the user
	// For this example, we will just simulate a successful signup
	if username == "" || password == "" || email == "" {
		return false, fmt.Errorf("Username, password, and email cannot be empty")
	}

	existingUser, errUser := repositories.ExisteUserByUsernameAndEmail(username, email)
	if existingUser != nil && errUser == nil {
		return false, fmt.Errorf("User with username or email already exists")
	} else if existingUser == nil && errUser != nil {
		return false, fmt.Errorf("Error checking existing user: %s", errUser.Error())
	}

	password = utils.HashPassword(password)

	newUser := models.User{
		Username:        username,
		Email:           email,
		PasswordHash:    password,
		Role:            "user", // Default role
		CreatedAt:       time.Now(),
		LastLogin:       nil,   // No last login yet
		IsBanned:        false, // Not banned by default
		BannedAt:        nil,   // No banned date yet
		ProfilePicture:  "",
		Biography:       "",
		MessageCount:    0, // Default message count
		DiscussionCount: 0, // Default discussion count
	}

	_, errNewUser := repositories.CreateUser(newUser) // Assuming this function exists in the repository
	if errNewUser != nil {
		return false, fmt.Errorf("Error creating new user: %s", errNewUser.Error())
	}

	return true, nil // Simulate a successful signup
}

func LoginService(password, username string) (string, error) {
	user, err := repositories.GetUserByEmailorUsername(username)
	if err != nil {
		return "", err
	}

	if !utils.ComparePasswords(user.PasswordHash, password) {
		return "", fmt.Errorf("invalid password")
	}

	token, err := utils.GenerateJWT(user.UserId)
	if err != nil {
		return "", err
	}
	fmt.Println("token ")

	return token, nil
}
