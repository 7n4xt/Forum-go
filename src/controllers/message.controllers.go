package controllers

import (
	"fmt"
	"forum-go/services"
	"forum-go/templates"
	"forum-go/utils"
	"net/http"
	"strconv"
	"strings"
)

// GetMessages handles the request to get all messages
func GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMessages handler called")

	// Check if user is authenticated (optional)
	user, _ := utils.GetUserFromRequest(r)

	// Get messages
	messages, err := services.GetAllMessagesService()
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Messages": messages,
		"User":     user,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "messages", data)
}

// GetMessageByID handles the request to get a single message
func GetMessageByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMessageByID handler called")

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Check if user is authenticated (optional)
	user, _ := utils.GetUserFromRequest(r)

	// Get message
	message, err := services.GetMessageByIDService(messageID)
	if err != nil {
		fmt.Println("Error fetching message:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if message == nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// Get discussion
	discussion, err := services.GetDiscussionByIDService(message.DiscussionId)
	if err != nil {
		fmt.Println("Error fetching discussion:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := map[string]interface{}{
		"Message":    message,
		"Discussion": discussion,
		"User":       user,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "message", data)
}

// CreateMessage handles the form submission to create a new message
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateMessage handler called")

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
	content := r.FormValue("content")
	discussionIdStr := r.FormValue("discussion_id")

	// Parse discussion ID
	discussionId, err := strconv.Atoi(discussionIdStr)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Handle file upload if present
	hasImage := false
	imagePath := ""
	file, handler, err := r.FormFile("image")
	if err == nil && file != nil {
		defer file.Close()
		// Process image upload (code not shown for brevity)
		// This would save the image to a directory and set hasImage and imagePath
		hasImage = true
		imagePath = "/public/images/" + handler.Filename
	}

	// Create message
	_, err = services.CreateMessageService(content, user.UserId, discussionId, hasImage, imagePath)
	if err != nil {
		http.Error(w, "Error creating message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussion page
	http.Redirect(w, r, "/discussions/"+discussionIdStr, http.StatusSeeOther)
}

// UpdateMessage handles the form submission to update a message
func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateMessage handler called")

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

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Get form values
	content := r.FormValue("content")

	// Update message
	err = services.UpdateMessageService(messageID, content, user.UserId)
	if err != nil {
		http.Error(w, "Error updating message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get message to determine which discussion to redirect to
	message, err := services.GetMessageByIDService(messageID)
	if err != nil {
		http.Error(w, "Error fetching message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussion page
	http.Redirect(w, r, "/discussions/"+strconv.Itoa(message.DiscussionId), http.StatusSeeOther)
}

// DeleteMessage handles the request to delete a message
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteMessage handler called")

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

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Get message to determine which discussion to redirect to
	message, err := services.GetMessageByIDService(messageID)
	if err != nil {
		http.Error(w, "Error fetching message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	discussionId := message.DiscussionId

	// Delete message
	isAdmin := user.Role == "admin"
	err = services.DeleteMessageService(messageID, user.UserId, isAdmin)
	if err != nil {
		http.Error(w, "Error deleting message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to discussion page
	http.Redirect(w, r, "/discussions/"+strconv.Itoa(discussionId), http.StatusSeeOther)
}
