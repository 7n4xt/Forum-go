package repositories

import (
	"database/sql"
	"fmt"
	"forum-go/config"
	"forum-go/models"
)

// CreateMessage inserts a new message into the database
func CreateMessage(message models.Message) (int, error) {
	query := "INSERT INTO `messages`(`content`, `author_id`, `discussion_id`, `has_image`, `image_path`) VALUES (?,?,?,?,?);"

	sqlResult, err := config.DbContext.Exec(
		query,
		message.Content,
		message.AuthorId,
		message.DiscussionId,
		message.HasImage,
		message.ImagePath,
	)

	if err != nil {
		return -1, fmt.Errorf("Erreur ajout message - Erreur : %s", err.Error())
	}

	id, err := sqlResult.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Erreur ajout message récupération ID - Erreur : %s", err.Error())
	}

	return int(id), nil
}

// ReadAllMessages retrieves all messages from the database
func ReadAllMessages() (*[]models.Message, error) {
	sqlResult, err := config.DbContext.Query("SELECT * FROM `messages`;")
	if err != nil {
		return nil, fmt.Errorf("Erreur récupération liste messages - Erreur : \n\t %s", err.Error())
	}
	defer sqlResult.Close()

	var listMessages []models.Message
	for sqlResult.Next() {
		var message models.Message

		err := sqlResult.Scan(
			&message.MessageId,
			&message.Content,
			&message.AuthorId,
			&message.DiscussionId,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.HasImage,
			&message.ImagePath,
			&message.LikeCount,
			&message.DislikeCount,
			&message.PopularityScore,
		)

		if err != nil {
			return nil, fmt.Errorf("Erreur récupération liste messages - Erreur : \n\t %s", err.Error())
		}

		listMessages = append(listMessages, message)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("Erreur récupération liste messages - Erreur : \n\t %s", sqlResult.Err())
	}

	return &listMessages, nil
}

// ReadMessageById retrieves a message by its ID
func ReadMessageById(messageId int) (*models.Message, error) {
	sqlResult := config.DbContext.QueryRow(
		"SELECT * FROM `messages` WHERE message_id = ?;",
		messageId,
	)

	var message models.Message

	err := sqlResult.Scan(
		&message.MessageId,
		&message.Content,
		&message.AuthorId,
		&message.DiscussionId,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.HasImage,
		&message.ImagePath,
		&message.LikeCount,
		&message.DislikeCount,
		&message.PopularityScore,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Erreur récupération message par ID - Erreur : \n\t %s", err.Error())
	}

	return &message, nil
}

// ReadMessagesByDiscussionId retrieves all messages from a discussion
func ReadMessagesByDiscussionId(discussionId int) (*[]models.Message, error) {
	sqlResult, err := config.DbContext.Query(
		"SELECT * FROM `messages` WHERE discussion_id = ? ORDER BY created_at DESC;",
		discussionId,
	)
	if err != nil {
		return nil, fmt.Errorf("Erreur récupération messages par discussion - Erreur : \n\t %s", err.Error())
	}
	defer sqlResult.Close()

	var listMessages []models.Message
	for sqlResult.Next() {
		var message models.Message

		err := sqlResult.Scan(
			&message.MessageId,
			&message.Content,
			&message.AuthorId,
			&message.DiscussionId,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.HasImage,
			&message.ImagePath,
			&message.LikeCount,
			&message.DislikeCount,
			&message.PopularityScore,
		)

		if err != nil {
			return nil, fmt.Errorf("Erreur récupération messages par discussion - Erreur : \n\t %s", err.Error())
		}

		listMessages = append(listMessages, message)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("Erreur récupération messages par discussion - Erreur : \n\t %s", sqlResult.Err())
	}

	return &listMessages, nil
}

// ReadMessagesByDiscussionIdWithSortAndPagination retrieves messages from a discussion with sorting and pagination
func ReadMessagesByDiscussionIdWithSortAndPagination(discussionId int, sortBy string, limit int, offset int) (*[]models.Message, error) {
	var query string

	// Build query based on sort type
	switch sortBy {
	case "popularity":
		query = "SELECT * FROM `messages` WHERE discussion_id = ? ORDER BY popularity_score DESC, created_at DESC"
	case "chronological":
		query = "SELECT * FROM `messages` WHERE discussion_id = ? ORDER BY created_at ASC"
	default: // default to newest first
		query = "SELECT * FROM `messages` WHERE discussion_id = ? ORDER BY created_at DESC"
	}

	// Add pagination
	if limit > 0 {
		query += " LIMIT ?"
		if offset > 0 {
			query += " OFFSET ?"
		}
	}

	// Prepare arguments
	var args []interface{}
	args = append(args, discussionId)
	if limit > 0 {
		args = append(args, limit)
		if offset > 0 {
			args = append(args, offset)
		}
	}

	sqlResult, err := config.DbContext.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Erreur récupération messages par discussion avec tri et pagination - Erreur : \n\t %s", err.Error())
	}
	defer sqlResult.Close()

	var listMessages []models.Message
	for sqlResult.Next() {
		var message models.Message

		err := sqlResult.Scan(
			&message.MessageId,
			&message.Content,
			&message.AuthorId,
			&message.DiscussionId,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.HasImage,
			&message.ImagePath,
			&message.LikeCount,
			&message.DislikeCount,
			&message.PopularityScore,
		)

		if err != nil {
			return nil, fmt.Errorf("Erreur récupération messages par discussion avec tri et pagination - Erreur : \n\t %s", err.Error())
		}

		listMessages = append(listMessages, message)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("Erreur récupération messages par discussion avec tri et pagination - Erreur : \n\t %s", sqlResult.Err())
	}

	return &listMessages, nil
}

// GetMessageCountByDiscussionId returns the total count of messages in a discussion
func GetMessageCountByDiscussionId(discussionId int) (int, error) {
	var count int
	err := config.DbContext.QueryRow("SELECT COUNT(*) FROM `messages` WHERE discussion_id = ?", discussionId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Erreur récupération nombre de messages - Erreur : \n\t %s", err.Error())
	}
	return count, nil
}

// ReadMessagesByAuthorId retrieves all messages by a specific author
func ReadMessagesByAuthorId(authorId int) (*[]models.Message, error) {
	sqlResult, err := config.DbContext.Query(
		"SELECT * FROM `messages` WHERE author_id = ? ORDER BY created_at DESC;",
		authorId,
	)
	if err != nil {
		return nil, fmt.Errorf("Erreur récupération messages par auteur - Erreur : \n\t %s", err.Error())
	}
	defer sqlResult.Close()

	var listMessages []models.Message
	for sqlResult.Next() {
		var message models.Message

		err := sqlResult.Scan(
			&message.MessageId,
			&message.Content,
			&message.AuthorId,
			&message.DiscussionId,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.HasImage,
			&message.ImagePath,
			&message.LikeCount,
			&message.DislikeCount,
			&message.PopularityScore,
		)

		if err != nil {
			return nil, fmt.Errorf("Erreur récupération messages par auteur - Erreur : \n\t %s", err.Error())
		}

		listMessages = append(listMessages, message)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("Erreur récupération messages par auteur - Erreur : \n\t %s", sqlResult.Err())
	}

	return &listMessages, nil
}

// UpdateMessage updates an existing message
func UpdateMessage(message models.Message) (bool, error) {
	query := "UPDATE `messages` SET `content` = ?, `has_image` = ?, `image_path` = ? WHERE `message_id` = ?;"

	sqlResult, err := config.DbContext.Exec(
		query,
		message.Content,
		message.HasImage,
		message.ImagePath,
		message.MessageId,
	)

	if err != nil {
		return false, fmt.Errorf("Erreur modification message %d - Erreur : \n\t %s", message.MessageId, err)
	}

	nbrRows, err := sqlResult.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erreur modification message - Erreur : \n\t %s", err)
	}

	if nbrRows == 0 {
		return false, nil
	}

	return true, nil
}

// DeleteMessage removes a message from the database and all related data
func DeleteMessage(messageId int) (bool, error) {
	// Begin transaction to ensure all related data is deleted atomically
	db := config.DbContext
	tx, err := db.Begin()
	if err != nil {
		return false, fmt.Errorf("Erreur démarrage transaction - Erreur : \n\t %s", err.Error())
	}
	defer tx.Rollback()

	// First delete any message reactions for this message
	_, err = tx.Exec("DELETE FROM message_reactions WHERE message_id = ?", messageId)
	if err != nil {
		return false, fmt.Errorf("Erreur suppression reactions message - Erreur : \n\t %s", err.Error())
	}

	// Delete any images associated with this message
	_, err = tx.Exec("DELETE FROM images WHERE message_id = ?", messageId)
	if err != nil {
		return false, fmt.Errorf("Erreur suppression images message - Erreur : \n\t %s", err.Error())
	}

	// Finally delete the message itself
	result, err := tx.Exec("DELETE FROM messages WHERE message_id = ?", messageId)
	if err != nil {
		return false, fmt.Errorf("Erreur suppression message - Erreur : \n\t %s", err.Error())
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Erreur suppression message - Erreur : \n\t %s", err.Error())
	}

	if rowsDeleted == 0 {
		return false, nil
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return false, fmt.Errorf("Erreur finalisation transaction - Erreur : \n\t %s", err.Error())
	}

	return true, nil
}

// DeleteMessagesByDiscussionId removes all messages from a discussion
func DeleteMessagesByDiscussionId(discussionId int) (int64, error) {
	query := "DELETE FROM `messages` WHERE `discussion_id` = ?;"
	result, err := config.DbContext.Exec(
		query,
		discussionId,
	)

	if err != nil {
		return 0, fmt.Errorf("Erreur suppression messages par discussion - Erreur : \n\t %s", err.Error())
	}

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Erreur suppression messages par discussion - Erreur : \n\t %s", err.Error())
	}

	return rowsDeleted, nil
}

// IncrementLikeCount increments the like count of a message
func IncrementLikeCount(messageId int) error {
	query := "UPDATE `messages` SET `like_count` = `like_count` + 1, `popularity_score` = `popularity_score` + 1 WHERE `message_id` = ?;"

	_, err := config.DbContext.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("Erreur incrémentation compteur de likes - Erreur : \n\t %s", err)
	}

	return nil
}

// IncrementDislikeCount increments the dislike count of a message
func IncrementDislikeCount(messageId int) error {
	query := "UPDATE `messages` SET `dislike_count` = `dislike_count` + 1, `popularity_score` = `popularity_score` - 1 WHERE `message_id` = ?;"

	_, err := config.DbContext.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("Erreur incrémentation compteur de dislikes - Erreur : \n\t %s", err)
	}

	return nil
}

// DecrementLikeCount decrements the like count of a message
func DecrementLikeCount(messageId int) error {
	query := "UPDATE `messages` SET `like_count` = `like_count` - 1, `popularity_score` = `popularity_score` - 1 WHERE `message_id` = ? AND `like_count` > 0;"

	_, err := config.DbContext.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("Erreur décrémentation compteur de likes - Erreur : \n\t %s", err)
	}

	return nil
}

// DecrementDislikeCount decrements the dislike count of a message
func DecrementDislikeCount(messageId int) error {
	query := "UPDATE `messages` SET `dislike_count` = `dislike_count` - 1, `popularity_score` = `popularity_score` + 1 WHERE `message_id` = ? AND `dislike_count` > 0;"

	_, err := config.DbContext.Exec(query, messageId)
	if err != nil {
		return fmt.Errorf("Erreur décrémentation compteur de dislikes - Erreur : \n\t %s", err)
	}

	return nil
}

// GetPopularMessages retrieves popular messages
func GetPopularMessages(limit int) (*[]models.Message, error) {
	sqlResult, err := config.DbContext.Query(
		"SELECT * FROM `messages` ORDER BY popularity_score DESC LIMIT ?;",
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("Erreur récupération messages populaires - Erreur : \n\t %s", err.Error())
	}
	defer sqlResult.Close()

	var listMessages []models.Message
	for sqlResult.Next() {
		var message models.Message

		err := sqlResult.Scan(
			&message.MessageId,
			&message.Content,
			&message.AuthorId,
			&message.DiscussionId,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.HasImage,
			&message.ImagePath,
			&message.LikeCount,
			&message.DislikeCount,
			&message.PopularityScore,
		)

		if err != nil {
			return nil, fmt.Errorf("Erreur récupération messages populaires - Erreur : \n\t %s", err.Error())
		}

		listMessages = append(listMessages, message)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("Erreur récupération messages populaires - Erreur : \n\t %s", sqlResult.Err())
	}

	return &listMessages, nil
}

// GetTotalMessagesCount returns the total number of messages
func GetTotalMessagesCount() (int, error) {
	var count int
	err := config.DbContext.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count)
	return count, err
}

// GetRecentMessages retrieves recent messages for admin review
func GetRecentMessages(limit int) (*[]models.Message, error) {
	sqlResult, err := config.DbContext.Query(
		"SELECT * FROM `messages` ORDER BY created_at DESC LIMIT ?;",
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recent messages: %s", err.Error())
	}
	defer sqlResult.Close()

	var listMessages []models.Message
	for sqlResult.Next() {
		var message models.Message

		err := sqlResult.Scan(
			&message.MessageId,
			&message.Content,
			&message.AuthorId,
			&message.DiscussionId,
			&message.CreatedAt,
			&message.UpdatedAt,
			&message.HasImage,
			&message.ImagePath,
			&message.LikeCount,
			&message.DislikeCount,
			&message.PopularityScore,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning recent messages: %s", err.Error())
		}

		listMessages = append(listMessages, message)
	}

	if sqlResult.Err() != nil {
		return nil, fmt.Errorf("error iterating recent messages: %s", sqlResult.Err())
	}

	return &listMessages, nil
}
