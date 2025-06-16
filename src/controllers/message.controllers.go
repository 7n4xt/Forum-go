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

	// Check if user is banned
	if user.IsBanned {
		http.Error(w, "Your account has been banned. You cannot post messages.", http.StatusForbidden)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Debug all form values
	fmt.Println("DEBUG - All form values:")
	for key, values := range r.Form {
		fmt.Printf("DEBUG - Form field %s: %v\n", key, values)
	}

	// Get form values
	content := r.FormValue("content")
	fmt.Printf("DEBUG - Content from form: '%s'\n", content)

	// Get discussion ID from URL path
	path := r.URL.Path
	fmt.Printf("DEBUG - URL path: '%s'\n", path)

	// Extract discussion ID from URL path
	parts := strings.Split(path, "/")
	var discussionIdStr string
	if len(parts) >= 3 && parts[len(parts)-2] == "create" {
		discussionIdStr = parts[len(parts)-1]
		fmt.Printf("DEBUG - Discussion ID from URL path: '%s'\n", discussionIdStr)
	} else {
		// Fallback to form value
		discussionIdStr = r.FormValue("discussion_id")
		fmt.Printf("DEBUG - Discussion ID from form: '%s'\n", discussionIdStr)
	}

	// Parse discussion ID
	discussionId, err := strconv.Atoi(discussionIdStr)
	if err != nil {
		fmt.Printf("DEBUG - Error parsing discussion ID: %v\n", err)
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	fmt.Printf("DEBUG - Parsed discussion ID: %d\n", discussionId)

	// Create message
	_, err = services.CreateMessageService(content, user.UserId, discussionId, false, "")
	if err != nil {
		fmt.Printf("DEBUG - Error creating message: %v\n", err)
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
	messageID, err := strconv.Atoi(pathParts[len(pathParts)-1])
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
	fmt.Println("URL path parts:", pathParts)
	if len(pathParts) < 3 {
		http.Error(w, "Invalid message ID format in URL", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid message ID: %s - Error: %v", pathParts[len(pathParts)-1], err), http.StatusBadRequest)
		return
	}

	// Get message to determine which discussion to redirect to
	message, err := services.GetMessageByIDService(messageID)
	if err != nil {
		http.Error(w, "Error fetching message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if message == nil {
		http.Error(w, "Message not found", http.StatusNotFound)
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
