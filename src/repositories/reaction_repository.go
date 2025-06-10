package repositories

import (
	"database/sql"
	"fmt"
	"forum-go/config"
)

// ReactionType enum
const (
	ReactionLike    = "like"
	ReactionDislike = "dislike"
)

// CreateReaction adds a new reaction (like or dislike) to a message
func CreateReaction(userId int, messageId int, reactionType string) error {
	// Begin transaction
	tx, err := config.DbContext.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First check if the user already has a reaction on this message
	var existingReactionType string
	var existingReactionId int
	err = tx.QueryRow(
		"SELECT reaction_id, reaction_type FROM message_reactions WHERE user_id = ? AND message_id = ?",
		userId, messageId,
	).Scan(&existingReactionId, &existingReactionType)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing reaction: %v", err)
	}

	// Handle existing reaction
	if err == nil {
		// User already has a reaction
		if existingReactionType == reactionType {
			// User is trying to add the same reaction they already have, so remove it
			_, err := tx.Exec(
				"DELETE FROM message_reactions WHERE reaction_id = ?",
				existingReactionId,
			)
			if err != nil {
				return fmt.Errorf("error removing existing reaction: %v", err)
			}

			// Update message counters
			if reactionType == ReactionLike {
				_, err = tx.Exec(
					"UPDATE messages SET like_count = like_count - 1, popularity_score = popularity_score - 1 WHERE message_id = ?",
					messageId,
				)
			} else {
				_, err = tx.Exec(
					"UPDATE messages SET dislike_count = dislike_count - 1, popularity_score = popularity_score + 1 WHERE message_id = ?",
					messageId,
				)
			}
			if err != nil {
				return fmt.Errorf("error updating message counts: %v", err)
			}
		} else {
			// User is changing their reaction (like -> dislike or dislike -> like)
			// Update the existing reaction
			_, err := tx.Exec(
				"UPDATE message_reactions SET reaction_type = ? WHERE reaction_id = ?",
				reactionType, existingReactionId,
			)
			if err != nil {
				return fmt.Errorf("error updating reaction: %v", err)
			}

			// Update message counters - we need to decrement the old reaction and increment the new one
			if reactionType == ReactionLike {
				// Changing from dislike to like
				_, err = tx.Exec(
					"UPDATE messages SET like_count = like_count + 1, dislike_count = dislike_count - 1, popularity_score = popularity_score + 2 WHERE message_id = ?",
					messageId,
				)
			} else {
				// Changing from like to dislike
				_, err = tx.Exec(
					"UPDATE messages SET like_count = like_count - 1, dislike_count = dislike_count + 1, popularity_score = popularity_score - 2 WHERE message_id = ?",
					messageId,
				)
			}
			if err != nil {
				return fmt.Errorf("error updating message counts: %v", err)
			}
		}
	} else {
		// No existing reaction, create a new one
		_, err = tx.Exec(
			"INSERT INTO message_reactions (user_id, message_id, reaction_type) VALUES (?, ?, ?)",
			userId, messageId, reactionType,
		)
		if err != nil {
			return fmt.Errorf("error creating reaction: %v", err)
		}

		// Update message counters
		if reactionType == ReactionLike {
			_, err = tx.Exec(
				"UPDATE messages SET like_count = like_count + 1, popularity_score = popularity_score + 1 WHERE message_id = ?",
				messageId,
			)
		} else {
			_, err = tx.Exec(
				"UPDATE messages SET dislike_count = dislike_count + 1, popularity_score = popularity_score - 1 WHERE message_id = ?",
				messageId,
			)
		}
		if err != nil {
			return fmt.Errorf("error updating message counts: %v", err)
		}
	}

	// Commit the transaction
	return tx.Commit()
}

// GetUserReaction gets a user's reaction to a message if it exists
func GetUserReaction(userId int, messageId int) (string, error) {
	var reactionType string
	err := config.DbContext.QueryRow(
		"SELECT reaction_type FROM message_reactions WHERE user_id = ? AND message_id = ?",
		userId, messageId,
	).Scan(&reactionType)

	if err == sql.ErrNoRows {
		return "", nil // No reaction
	}

	if err != nil {
		return "", fmt.Errorf("error getting user reaction: %v", err)
	}

	return reactionType, nil
}

// GetMessageReactions gets reactions stats for a message
func GetMessageReactions(messageId int) (int, int, error) {
	var likeCount, dislikeCount int
	err := config.DbContext.QueryRow(
		"SELECT like_count, dislike_count FROM messages WHERE message_id = ?",
		messageId,
	).Scan(&likeCount, &dislikeCount)

	if err != nil {
		return 0, 0, fmt.Errorf("error getting message reactions: %v", err)
	}

	return likeCount, dislikeCount, nil
}

// DeleteReactionsForMessage deletes all reactions for a message
func DeleteReactionsForMessage(messageId int) error {
	_, err := config.DbContext.Exec(
		"DELETE FROM message_reactions WHERE message_id = ?",
		messageId,
	)
	return err
}
