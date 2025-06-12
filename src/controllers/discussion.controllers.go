package controllers

import (
	"fmt"
	"forum-go/models"
	"forum-go/services"
	"forum-go/templates"
	"forum-go/utils"
	"net/http"
	"strconv"
	"strings"
)

// GetAllDiscussions handles the request to get all discussions
func GetAllDiscussions(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated (optional)
	user, _ := utils.GetUserFromRequest(r)

	// Get query parameters for filtering
	status := r.URL.Query().Get("status")
	categoryIdStr := r.URL.Query().Get("category")
	authorIdStr := r.URL.Query().Get("author")

	// Parse parameters
	categoryId := 0
	if categoryIdStr != "" {
		var err error
		categoryId, err = strconv.Atoi(categoryIdStr)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}
	}

	authorId := 0
	if authorIdStr != "" {
		var err error
		authorId, err = strconv.Atoi(authorIdStr)
		if err != nil {
			http.Error(w, "Invalid author ID", http.StatusBadRequest)
			return
		}
	}

	// Get discussions
	discussions, err := services.GetAllDiscussionsService(status, categoryId, authorId, 20, 0)
	if err != nil {
		http.Error(w, "Error fetching discussions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get categories for the filter dropdown
	categories, err := services.GetAllCategoriesService()
	if err != nil {
		fmt.Println("Error fetching categories:", err)
		// Continue without categories
	}

	// Prepare data for template
	data := map[string]any{
		"Discussions": discussions,
		"Categories":  categories,
		"User":        user,
		"Status":      status,
		"CategoryId":  categoryId,
		"AuthorId":    authorId,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "discussions", data)
}

// GetDiscussionByID handles the request to get a single discussion by ID
func GetDiscussionByID(w http.ResponseWriter, r *http.Request) {
	// Extract discussion ID from URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	discussionIdStr := parts[2]
	fmt.Printf("DEBUG - Discussion ID from URL: '%s'\n", discussionIdStr)

	discussionId, err := strconv.Atoi(discussionIdStr)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Check if user is authenticated (optional)
	user, _ := utils.GetUserFromRequest(r)

	// Get discussion
	discussion, err := services.GetDiscussionByIDService(discussionId)
	if err != nil {
		fmt.Printf("DEBUG - Error fetching discussion: %v\n", err)
		http.Error(w, "Error fetching discussion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("DEBUG - Discussion object: %+v\n", discussion)

	// Get messages for this discussion
	messages, err := services.GetMessagesByDiscussionIdService(discussionId)
	if err != nil {
		http.Error(w, "Error fetching messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add user reactions if user is logged in
	var messagesWithReactions []models.MessageWithReaction
	var userId int = 0
	if user != nil {
		userId = user.UserId
	}

	messagesWithReactions, err = services.EnrichMessagesWithUserReactions(messages, userId)
	if err != nil {
		http.Error(w, "Error getting reactions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Discussion": discussion,
		"Messages":   messagesWithReactions,
		"User":       user,
		"CanPost":    user != nil && discussion.Status == "open",
	}

	fmt.Printf("DEBUG - Discussion ID in template data: %d\n", discussion.DiscussionId)

	// Render template
	templates.Temp.ExecuteTemplate(w, "discussion", data)
}

// CreateDiscussionForm handles the form submission to create a new discussion
func CreateDiscussionForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/discussions", http.StatusSeeOther)
		return
	}

	// Check if user is authenticated
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	categoryIdsStr := r.Form["categories"]

	// Convert category IDs to integers
	var categoryIds []int
	for _, idStr := range categoryIdsStr {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid category ID: "+idStr, http.StatusBadRequest)
			return
		}
		categoryIds = append(categoryIds, id)
	}

	// Create discussion
	_, err = services.CreateDiscussionService(title, description, categoryIds, user.UserId)
	if err != nil {
		http.Error(w, "Error creating discussion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussions list
	http.Redirect(w, r, "/discussions", http.StatusSeeOther)
}

// NewDiscussionPage renders the page to create a new discussion
func NewDiscussionPage(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get categories for the form
	categories, err := services.GetAllCategoriesService()
	if err != nil {
		http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Categories": categories,
		"User":       user,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "new_discussion", data)
}

// UpdateDiscussionContent handles the request to update a discussion's content
func UpdateDiscussionContent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateDiscussionContent handler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	discussionIdStr := r.FormValue("discussion_id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	// Validate inputs
	if title == "" || description == "" {
		http.Error(w, "Title and description are required", http.StatusBadRequest)
		return
	}

	// Convert discussion ID to integer
	discussionId, err := strconv.Atoi(discussionIdStr)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Check if user is admin
	isAdmin := user.Role == "admin"

	// Update discussion content
	err = services.UpdateDiscussionContentService(discussionId, title, description, user.UserId, isAdmin)
	if err != nil {
		http.Error(w, "Error updating discussion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussion page
	http.Redirect(w, r, "/discussions/"+discussionIdStr, http.StatusSeeOther)
}

// UpdateDiscussionStatus handles the request to update a discussion's status
func UpdateDiscussionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get form values
	discussionIdStr := r.FormValue("discussion_id")
	status := r.FormValue("status")

	// Convert discussion ID to integer
	discussionId, err := strconv.Atoi(discussionIdStr)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Check if user is admin
	isAdmin := user.Role == "admin"

	// Update discussion status
	err = services.UpdateDiscussionStatusService(discussionId, status, user.UserId, isAdmin)
	if err != nil {
		http.Error(w, "Error updating discussion status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussion page
	http.Redirect(w, r, "/discussions/"+discussionIdStr, http.StatusSeeOther)
}

// DeleteDiscussion handles the request to delete a discussion
func DeleteDiscussion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteDiscussion handler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Extract ID from URL path or form
	var discussionId int
	if r.URL.Path == "/discussions/delete" {
		// Get from form
		discussionIdStr := r.FormValue("discussion_id")
		discussionId, err = strconv.Atoi(discussionIdStr)
		if err != nil {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}
	} else {
		// Get from URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}
		discussionId, err = strconv.Atoi(pathParts[2])
		if err != nil {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}
	}

	// Check if user is admin
	isAdmin := user.Role == "admin"

	// Delete discussion
	err = services.DeleteDiscussionService(discussionId, user.UserId, isAdmin)
	if err != nil {
		http.Error(w, "Error deleting discussion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussions page
	http.Redirect(w, r, "/discussions", http.StatusSeeOther)
}
