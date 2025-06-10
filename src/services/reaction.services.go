package services

import (
	"errors"
	"fmt"
	"forum-go/models"
	"forum-go/repositories"
)

// ReactToMessageService handles message reactions (likes/dislikes)
func ReactToMessageService(userId int, messageId int, reactionType string) error {
	// Validate reaction type
	if reactionType != repositories.ReactionLike && reactionType != repositories.ReactionDislike {
		return errors.New("invalid reaction type")
	}

	// Check if message exists
	message, err := repositories.ReadMessageById(messageId)
	if err != nil {
		return fmt.Errorf("error checking message: %v", err)
	}
	if message == nil {
		return errors.New("message not found")
	}

	// Create reaction
	return repositories.CreateReaction(userId, messageId, reactionType)
}

// GetUserReactionToMessageService gets a user's reaction to a message
func GetUserReactionToMessageService(userId int, messageId int) (string, error) {
	return repositories.GetUserReaction(userId, messageId)
}

// GetMessageReactionsService gets reactions for a message
func GetMessageReactionsService(messageId int) (int, int, error) {
	return repositories.GetMessageReactions(messageId)
}

// EnrichMessagesWithUserReactions adds user's reaction to a list of messages
func EnrichMessagesWithUserReactions(messages []models.Message, userId int) ([]models.MessageWithReaction, error) {
	if userId == 0 {
		// Convert regular messages to MessageWithReaction with no reactions
		result := make([]models.MessageWithReaction, len(messages))
		for i, msg := range messages {
			result[i] = models.MessageWithReaction{
				Message:      msg,
				UserReaction: "",
			}
		}
		return result, nil
	}

	// Get user reactions for each message
	result := make([]models.MessageWithReaction, len(messages))
	for i, msg := range messages {
		// Get user's reaction to this message
		reaction, err := GetUserReactionToMessageService(userId, msg.MessageId)
		if err != nil {
			return nil, err
		}

		result[i] = models.MessageWithReaction{
			Message:      msg,
			UserReaction: reaction,
		}
	}

	return result, nil
}
