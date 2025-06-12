package services

import (
	"errors"
	"forum-go/models"
	"forum-go/repositories"
	"time"
)

// CreateDiscussionService creates a new discussion
func CreateDiscussionService(title, description string, categoryIds []int, authorId int) (int64, error) {
	// Validate inputs
	if title == "" {
		return 0, errors.New("title cannot be empty")
	}
	if description == "" {
		return 0, errors.New("description cannot be empty")
	}
	if len(categoryIds) == 0 {
		return 0, errors.New("at least one category must be selected")
	}

	// Create discussion object
	discussion := models.Discussion{
		Title:       title,
		Description: description,
		Status:      "open",
		Visibility:  "public",
		AuthorId:    authorId,
		CategoryIds: categoryIds,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Call repository to save discussion
	return repositories.CreateDiscussion(discussion)
}

// GetDiscussionByIDService retrieves a discussion by its ID
func GetDiscussionByIDService(discussionId int) (models.Discussion, error) {
	return repositories.GetDiscussionByID(discussionId)
}

// GetAllDiscussionsService retrieves all discussions with optional filtering
func GetAllDiscussionsService(status string, categoryId int, authorId int, limit int, offset int) ([]models.Discussion, error) {
	return repositories.GetAllDiscussions(status, categoryId, authorId, limit, offset)
}

// UpdateDiscussionStatusService updates the status of a discussion
func UpdateDiscussionStatusService(discussionId int, status string, userId int, isAdmin bool) error {
	// Get discussion to check ownership
	discussion, err := repositories.GetDiscussionByID(discussionId)
	if err != nil {
		return err
	}

	// Check if user is discussion owner or admin
	if discussion.AuthorId != userId && !isAdmin {
		return errors.New("only the discussion owner or an admin can update the status")
	}

	return repositories.UpdateDiscussionStatus(discussionId, status, userId)
}

// UpdateDiscussionContentService updates the title and description of a discussion
func UpdateDiscussionContentService(discussionId int, title, description string, userId int, isAdmin bool) error {
	// Validate inputs
	if title == "" {
		return errors.New("title cannot be empty")
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}

	// Get discussion to check ownership
	discussion, err := repositories.GetDiscussionByID(discussionId)
	if err != nil {
		return err
	}

	// Check if user is discussion owner or admin
	if discussion.AuthorId != userId && !isAdmin {
		return errors.New("only the discussion owner or an admin can update the discussion")
	}

	// Check if discussion is still open
	if discussion.Status != "open" {
		return errors.New("cannot edit a closed or archived discussion")
	}

	// Update discussion
	return repositories.UpdateDiscussionContent(discussionId, title, description)
}

// DeleteDiscussionService deletes a discussion
func DeleteDiscussionService(discussionId int, userId int, isAdmin bool) error {
	// Get discussion to check ownership
	discussion, err := repositories.GetDiscussionByID(discussionId)
	if err != nil {
		return err
	}

	// Check if user is discussion owner or admin
	if discussion.AuthorId != userId && !isAdmin {
		return errors.New("only the discussion owner or an admin can delete the discussion")
	}

	// First delete all messages in the discussion
	_, err = repositories.DeleteMessagesByDiscussionId(discussionId)
	if err != nil {
		return err
	}

	// Then delete the discussion
	return repositories.DeleteDiscussion(discussionId)
}

// GetAllCategoriesService retrieves all categories
func GetAllCategoriesService() ([]models.Category, error) {
	return repositories.GetAllCategories()
}
