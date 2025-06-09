package models

import (
	"time"
)

// User represents a user in the database
type User struct {
	UserId          int        // Primary key
	Username        string     // Unique username
	Email           string     // Unique email
	PasswordHash    string     // SHA512 hash of password
	Role            string     // "user" or "admin"
	CreatedAt       time.Time  // Account creation timestamp
	LastLogin       *time.Time // Pointer for nullable last_login field
	IsBanned        bool       // Whether user is banned
	BannedAt        *time.Time // Pointer for nullable banned_at field
	ProfilePicture  string     // URL or path to profile picture
	Biography       string     // User's biography text
	MessageCount    int        // Count of messages sent by user
	DiscussionCount int        // Count of discussions started by user
}

// Message represents a message in a discussion
type Message struct {
	MessageId       int       // Primary key
	Content         string    // Message text content
	AuthorId        int       // Foreign key to User
	Author          *User     // Author information (populated when needed)
	DiscussionId    int       // Foreign key to Discussion
	CreatedAt       time.Time // Message creation timestamp
	UpdatedAt       time.Time // Last update timestamp
	HasImage        bool      // Whether message has an attached image
	ImagePath       string    // Path to attached image (if any)
	LikeCount       int       // Number of likes
	DislikeCount    int       // Number of dislikes
	PopularityScore int       // like_count - dislike_count
}

// Discussion represents a discussion thread
type Discussion struct {
	DiscussionId int       // Primary key
	Title        string    // Discussion title
	Description  string    // Discussion description
	Status       string    // "open", "closed", or "archived"
	Visibility   string    // "public" or "private"
	AuthorId     int       // Foreign key to User
	Author       *User     // Author information (populated when needed)
	CreatedAt    time.Time // Discussion creation timestamp
	UpdatedAt    time.Time // Last update timestamp
	ViewCount    int       // Number of views
	MessageCount int       // Number of messages in the discussion
	Categories   []string  // Categories/tags for the discussion
	CategoryIds  []int     // Category IDs (for database operations)
}

// Category represents a discussion category/tag
type Category struct {
	CategoryId  int       // Primary key
	Name        string    // Category name
	Description string    // Category description
	Color       string    // Hex color code
	CreatedAt   time.Time // Category creation timestamp
}
