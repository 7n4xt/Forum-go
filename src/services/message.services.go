package services

import (
	"errors"
	"forum-go/models"
	"forum-go/repositories"
)

// GetAllMessagesService returns all messages
func GetAllMessagesService() ([]models.Message, error) {
	messages, err := repositories.ReadAllMessages()
	if err != nil {
		return nil, err
	}

	// If no messages, return empty slice instead of nil
	if messages == nil {
		return []models.Message{}, nil
	}

	return *messages, nil
}

// GetMessageByIDService returns a message by its ID
func GetMessageByIDService(messageID int) (*models.Message, error) {
	return repositories.ReadMessageById(messageID)
}

// GetMessagesByDiscussionIdService returns all messages in a discussion
func GetMessagesByDiscussionIdService(discussionID int) ([]models.Message, error) {
	messages, err := repositories.ReadMessagesByDiscussionId(discussionID)
	if err != nil {
		return nil, err
	}

	// If no messages, return empty slice instead of nil
	if messages == nil {
		return []models.Message{}, nil
	}

	// Get author information for each message
	result := *messages
	for i := range result {
		author, err := repositories.ReadUserById(result[i].AuthorId)
		if err == nil && author != nil {
			result[i].Author = author
		}
	}

	return result, nil
}

// CreateMessageService creates a new message
func CreateMessageService(content string, authorID int, discussionID int, hasImage bool, imagePath string) (int, error) {
	if content == "" {
		return 0, errors.New("content is required")
	}

	// Check if discussion exists and is open
	discussion, err := repositories.GetDiscussionByID(discussionID)
	if err != nil {
		return 0, err
	}

	if discussion.Status != "open" {
		return 0, errors.New("cannot post in a closed or archived discussion")
	}

	// Create message
	message := models.Message{
		Content:      content,
		AuthorId:     authorID,
		DiscussionId: discussionID,
		HasImage:     hasImage,
		ImagePath:    imagePath,
	}

	return repositories.CreateMessage(message)
}

// UpdateMessageService updates an existing message
func UpdateMessageService(messageID int, content string, userID int) error {
	// Get existing message
	message, err := repositories.ReadMessageById(messageID)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.New("message not found")
	}

	// Check if user is the author
	if message.AuthorId != userID {
		return errors.New("you can only edit your own messages")
	}

	// Check if discussion is still open
	discussion, err := repositories.GetDiscussionByID(message.DiscussionId)
	if err != nil {
		return err
	}

	if discussion.Status != "open" {
		return errors.New("cannot edit messages in a closed or archived discussion")
	}

	// Update message
	message.Content = content

	_, err = repositories.UpdateMessage(*message)
	return err
}

// DeleteMessageService deletes a message
func DeleteMessageService(messageID int, userID int, isAdmin bool) error {
	// Get existing message
	message, err := repositories.ReadMessageById(messageID)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.New("message not found")
	}

	// Check if user is the author or an admin
	if message.AuthorId != userID && !isAdmin {
		return errors.New("you can only delete your own messages")
	}

	// Delete message
	_, err = repositories.DeleteMessage(messageID)
	return err
}
