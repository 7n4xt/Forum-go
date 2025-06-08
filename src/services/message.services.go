package services

import (
	"errors"
	"fmt"
)

// Message represents a forum message
type Message struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// In-memory storage for messages (replace with database in production)
var messages = make(map[int]Message)
var lastID = 0

// GetAllMessagesService returns all messages
func GetAllMessagesService() ([]Message, error) {
	var result []Message
	for _, msg := range messages {
		result = append(result, msg)
	}
	return result, nil
}

// GetMessageByIDService returns a message by its ID
func GetMessageByIDService(messageID int) (Message, error) {
	message, exists := messages[messageID]
	if !exists {
		return Message{}, errors.New("message not found")
	}
	return message, nil
}

// CreateMessageService creates a new message
func CreateMessageService(userID int, title, content string) (int, error) {
	if title == "" || content == "" {
		return 0, errors.New("title and content are required")
	}

	lastID++
	newMessage := Message{
		ID:      lastID,
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	messages[lastID] = newMessage
	fmt.Printf("Created message with ID: %d\n", lastID)
	return lastID, nil
}

// UpdateMessageService updates an existing message
func UpdateMessageService(messageID int, title, content string) error {
	message, exists := messages[messageID]
	if !exists {
		return errors.New("message not found")
	}

	if title != "" {
		message.Title = title
	}

	if content != "" {
		message.Content = content
	}

	messages[messageID] = message
	return nil
}

// DeleteMessageService deletes a message
func DeleteMessageService(messageID int) error {
	if _, exists := messages[messageID]; !exists {
		return errors.New("message not found")
	}

	delete(messages, messageID)
	return nil
}
