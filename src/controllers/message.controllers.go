package controllers

import (
	"encoding/json"
	"fmt"
	"forum-go/services"
	"net/http"
	"strconv"
	"strings"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMessages handler called")

	messages, err := services.GetAllMessagesService()
	if err != nil {
		fmt.Println("Error fetching messages:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func GetMessageByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMessageByID handler called")

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) == 0 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	message, err := services.GetMessageByIDService(messageID)
	if err != nil {
		fmt.Println("Error fetching message:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateMessage handler called")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body or form values depending on content type
	var content, title string
	userID := 0 // This would typically come from the authenticated user's token

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		content = fmt.Sprintf("%v", requestBody["content"])
		title = fmt.Sprintf("%v", requestBody["title"])
		if userIDVal, ok := requestBody["user_id"]; ok {
			userID, _ = strconv.Atoi(fmt.Sprintf("%v", userIDVal))
		}
	} else {
		// Form data
		r.ParseForm()
		content = r.FormValue("content")
		title = r.FormValue("title")
		userIDStr := r.FormValue("user_id")
		if userIDStr != "" {
			userID, _ = strconv.Atoi(userIDStr)
		}
	}

	messageID, err := services.CreateMessageService(userID, title, content)
	if err != nil {
		fmt.Println("Error creating message:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Message created successfully",
		"id":      messageID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateMessage handler called")

	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) == 0 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	var content, title string

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		if val, ok := requestBody["content"]; ok {
			content = fmt.Sprintf("%v", val)
		}
		if val, ok := requestBody["title"]; ok {
			title = fmt.Sprintf("%v", val)
		}
	} else {
		// Form data
		r.ParseForm()
		content = r.FormValue("content")
		title = r.FormValue("title")
	}

	err = services.UpdateMessageService(messageID, title, content)
	if err != nil {
		fmt.Println("Error updating message:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Message updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DeleteMessage handler called")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) == 0 {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	err = services.DeleteMessageService(messageID)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Message deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
