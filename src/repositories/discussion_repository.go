package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"forum-go/config"
	"forum-go/models"
	"strings"
	"time"
)

// CreateDiscussion inserts a new discussion into the database
func CreateDiscussion(discussion models.Discussion) (int64, error) {
	db := config.DbContext

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Insert discussion
	result, err := tx.Exec(
		"INSERT INTO discussions (title, description, status, visibility, author_id) VALUES (?, ?, ?, ?, ?)",
		discussion.Title, discussion.Description, discussion.Status, discussion.Visibility, discussion.AuthorId,
	)
	if err != nil {
		return 0, err
	}

	// Get the inserted discussion ID
	discussionId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert category associations if any
	if len(discussion.CategoryIds) > 0 {
		// Prepare statement for inserting categories
		stmt, err := tx.Prepare("INSERT INTO discussion_categories (discussion_id, category_id) VALUES (?, ?)")
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		// Insert each category association
		for _, categoryId := range discussion.CategoryIds {
			_, err = stmt.Exec(discussionId, categoryId)
			if err != nil {
				return 0, err
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return discussionId, nil
}

// GetDiscussionByID retrieves a discussion by its ID
func GetDiscussionByID(discussionId int) (models.Discussion, error) {
	db := config.DbContext
	var discussion models.Discussion

	// Query for discussion details
	err := db.QueryRow(`
		SELECT d.discussion_id, d.title, d.description, d.status, d.visibility, 
		       d.author_id, d.created_at, d.updated_at, d.view_count, d.message_count
		FROM discussions d
		WHERE d.discussion_id = ?
	`, discussionId).Scan(
		&discussion.DiscussionId, &discussion.Title, &discussion.Description, &discussion.Status, &discussion.Visibility,
		&discussion.AuthorId, &discussion.CreatedAt, &discussion.UpdatedAt, &discussion.ViewCount, &discussion.MessageCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return discussion, errors.New("discussion not found")
		}
		return discussion, err
	}

	// Get author information
	author, err := ReadUserById(discussion.AuthorId)
	if err == nil && author != nil {
		discussion.Author = author
	}

	// Get categories for the discussion
	rows, err := db.Query(`
		SELECT c.category_id, c.name
		FROM categories c
		JOIN discussion_categories dc ON c.category_id = dc.category_id
		WHERE dc.discussion_id = ?
	`, discussionId)
	if err != nil {
		return discussion, err
	}
	defer rows.Close()

	var categories []string
	var categoryIds []int
	for rows.Next() {
		var categoryId int
		var categoryName string
		if err := rows.Scan(&categoryId, &categoryName); err != nil {
			return discussion, err
		}
		categories = append(categories, categoryName)
		categoryIds = append(categoryIds, categoryId)
	}
	discussion.Categories = categories
	discussion.CategoryIds = categoryIds

	// Increment view count
	_, err = db.Exec("UPDATE discussions SET view_count = view_count + 1 WHERE discussion_id = ?", discussionId)
	if err != nil {
		// Log error but continue, as this is not critical
		fmt.Printf("Error incrementing view count: %v\n", err)
	}

	return discussion, nil
}

// GetAllDiscussions retrieves all discussions with optional filtering
func GetAllDiscussions(status string, categoryId int, authorId int, limit int, offset int) ([]models.Discussion, error) {
	db := config.DbContext
	var discussions []models.Discussion

	// Build query with filters
	query := `
		SELECT DISTINCT d.discussion_id, d.title, d.description, d.status, d.visibility, 
		       d.author_id, d.created_at, d.updated_at, d.view_count, d.message_count
		FROM discussions d
	`

	// Add joins if filtering by category
	if categoryId > 0 {
		query += `JOIN discussion_categories dc ON d.discussion_id = dc.discussion_id `
	}

	// Add WHERE clauses
	whereConditions := []string{}
	args := []interface{}{}

	// Filter by status
	if status != "" {
		whereConditions = append(whereConditions, "d.status = ?")
		args = append(args, status)
	} else {
		// By default, don't show archived discussions
		whereConditions = append(whereConditions, "d.status != 'archived'")
	}

	// Filter by visibility (always show public discussions)
	whereConditions = append(whereConditions, "d.visibility = 'public'")

	// Filter by category
	if categoryId > 0 {
		whereConditions = append(whereConditions, "dc.category_id = ?")
		args = append(args, categoryId)
	}

	// Filter by author
	if authorId > 0 {
		whereConditions = append(whereConditions, "d.author_id = ?")
		args = append(args, authorId)
	}

	// Add WHERE clause if conditions exist
	if len(whereConditions) > 0 {
		query += "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Add ordering and limits
	query += " ORDER BY d.updated_at DESC"
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
		if offset > 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process results
	for rows.Next() {
		var discussion models.Discussion
		err := rows.Scan(
			&discussion.DiscussionId, &discussion.Title, &discussion.Description, &discussion.Status, &discussion.Visibility,
			&discussion.AuthorId, &discussion.CreatedAt, &discussion.UpdatedAt, &discussion.ViewCount, &discussion.MessageCount,
		)
		if err != nil {
			return nil, err
		}

		// Get author information
		author, err := ReadUserById(discussion.AuthorId)
		if err == nil && author != nil {
			discussion.Author = author
		}

		// Get categories for the discussion
		categoryRows, err := db.Query(`
			SELECT c.category_id, c.name
			FROM categories c
			JOIN discussion_categories dc ON c.category_id = dc.category_id
			WHERE dc.discussion_id = ?
		`, discussion.DiscussionId)
		if err != nil {
			return nil, err
		}

		var categories []string
		var categoryIds []int
		for categoryRows.Next() {
			var categoryId int
			var categoryName string
			if err := categoryRows.Scan(&categoryId, &categoryName); err != nil {
				categoryRows.Close()
				return nil, err
			}
			categories = append(categories, categoryName)
			categoryIds = append(categoryIds, categoryId)
		}
		categoryRows.Close()
		discussion.Categories = categories
		discussion.CategoryIds = categoryIds

		discussions = append(discussions, discussion)
	}

	return discussions, nil
}

// UpdateDiscussionStatus updates the status of a discussion
func UpdateDiscussionStatus(discussionId int, status string, userId int) error {
	db := config.DbContext

	// First check if the user is the author or an admin
	var authorId int
	var userRole string

	// Get the discussion's author ID
	err := db.QueryRow("SELECT author_id FROM discussions WHERE discussion_id = ?", discussionId).Scan(&authorId)
	if err != nil {
		return err
	}

	// Get the user's role
	err = db.QueryRow("SELECT role FROM users WHERE user_id = ?", userId).Scan(&userRole)
	if err != nil {
		return err
	}

	// Check if user has permission
	if authorId != userId && userRole != "admin" {
		return errors.New("unauthorized: only the author or an admin can update discussion status")
	}

	// Validate status
	if status != "open" && status != "closed" && status != "archived" {
		return errors.New("invalid status: must be 'open', 'closed', or 'archived'")
	}

	// Update the discussion status
	_, err = db.Exec(
		"UPDATE discussions SET status = ?, updated_at = ? WHERE discussion_id = ?",
		status, time.Now(), discussionId,
	)
	return err
}

// GetAllCategories retrieves all categories
func GetAllCategories() ([]models.Category, error) {
	db := config.DbContext
	var categories []models.Category

	rows, err := db.Query("SELECT category_id, name, description, color, created_at FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.CategoryId, &category.Name, &category.Description, &category.Color, &category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// DeleteDiscussion removes a discussion from the database
func DeleteDiscussion(discussionId int) error {
	db := config.DbContext

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// First delete all category associations
	_, err = tx.Exec("DELETE FROM discussion_categories WHERE discussion_id = ?", discussionId)
	if err != nil {
		return err
	}

	// Then delete the discussion itself
	result, err := tx.Exec("DELETE FROM discussions WHERE discussion_id = ?", discussionId)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("discussion not found")
	}

	// Commit the transaction
	return tx.Commit()
}
