package controllers

import (
	"fmt"
	"forum-go/repositories"
	"forum-go/services"
	"forum-go/utils"
	"net/http"
	"strconv"
	"strings"
)

// ReactToMessage handles reactions (like/dislike) to messages
func ReactToMessage(w http.ResponseWriter, r *http.Request) {
	// Only accept POST method
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

	// Get path parts to extract message ID and reaction type
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Extract message ID from URL
	messageID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Extract reaction type from URL
	reactionType := pathParts[4]
	if reactionType != repositories.ReactionLike && reactionType != repositories.ReactionDislike {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	// Get the message to know which discussion to redirect to
	message, err := services.GetMessageByIDService(messageID)
	if err != nil || message == nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	// Process reaction
	err = services.ReactToMessageService(user.UserId, messageID, reactionType)
	if err != nil {
		http.Error(w, "Error processing reaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to the discussion page
	http.Redirect(w, r, fmt.Sprintf("/discussions/%d#message-%d", message.DiscussionId, messageID), http.StatusSeeOther)
}
