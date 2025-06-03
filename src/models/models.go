package models

import (
    "time"
)



// User represents a user in the database
type User struct {
    UserId           int       // Primary key
    Username         string    // Unique username
    Email            string    // Unique email
    PasswordHash     string    // SHA512 hash of password
    Role             string    // "user" or "admin"
    CreatedAt        time.Time // Account creation timestamp
    LastLogin        *time.Time // Pointer for nullable last_login field
    IsBanned         bool      // Whether user is banned
    BannedAt         *time.Time // Pointer for nullable banned_at field
    ProfilePicture   string    // URL or path to profile picture
    Biography        string    // User's biography text
    MessageCount     int       // Count of messages sent by user
    DiscussionCount  int       // Count of discussions started by user
}