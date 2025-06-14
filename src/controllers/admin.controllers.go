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

// AdminDashboard renders the admin dashboard page
func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Get dashboard statistics
	stats, err := services.GetAdminStatsService()
	if err != nil {
		http.Error(w, "Error fetching dashboard stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get recent discussions for admin review
	discussions, _, err := services.GetAllDiscussionsWithPaginationService("", 0, 0, 10, 0)
	if err != nil {
		fmt.Println("Error fetching discussions for admin:", err)
		discussions = []models.Discussion{}
	}

	// Get recent messages for admin review
	messages, err := services.GetRecentMessagesService(10)
	if err != nil {
		fmt.Println("Error fetching messages for admin:", err)
		messages = []models.Message{}
	}

	// Get recent users
	users, err := services.GetRecentUsersService(10)
	if err != nil {
		fmt.Println("Error fetching users for admin:", err)
		users = []models.User{}
	}

	// Prepare data for template
	data := map[string]interface{}{
		"User":        user,
		"Stats":       stats,
		"Discussions": discussions,
		"Messages":    messages,
		"Users":       users,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "admin_dashboard", data)
}

// AdminManageDiscussions renders the discussions management page
func AdminManageDiscussions(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Get query parameters for filtering and pagination
	status := r.URL.Query().Get("status")
	searchQuery := r.URL.Query().Get("search")
	pageStr := r.URL.Query().Get("page")

	// Parse pagination parameters
	limit := 20 // default for admin view
	page := 1
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit

	// Get discussions with search and pagination
	var discussions []models.Discussion
	var totalCount int

	if searchQuery != "" {
		discussions, totalCount, err = services.SearchDiscussionsService(searchQuery, status, 0, 0, limit, offset)
	} else {
		discussions, totalCount, err = services.GetAllDiscussionsWithPaginationService(status, 0, 0, limit, offset)
	}

	if err != nil {
		http.Error(w, "Error fetching discussions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate pagination info
	totalPages := (totalCount + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	// Prepare data for template
	data := map[string]interface{}{
		"User":        user,
		"Discussions": discussions,
		"Status":      status,
		"SearchQuery": searchQuery,
		"Page":        page,
		"TotalCount":  totalCount,
		"TotalPages":  totalPages,
		"HasNext":     hasNext,
		"HasPrev":     hasPrev,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "admin_discussions", data)
}

// AdminManageUsers renders the users management page
func AdminManageUsers(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Get query parameters for filtering and pagination
	searchQuery := r.URL.Query().Get("search")
	statusFilter := r.URL.Query().Get("status") // "active", "banned", or ""
	pageStr := r.URL.Query().Get("page")

	// Parse pagination parameters
	limit := 20 // default for admin view
	page := 1
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	offset := (page - 1) * limit

	// Get users with search and pagination
	users, totalCount, err := services.GetUsersForAdminService(searchQuery, statusFilter, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get stats for the users page
	stats, err := services.GetAdminStatsService()
	if err != nil {
		fmt.Println("Error fetching stats for admin users page:", err)
		// Continue without stats
	}

	// Calculate pagination info
	totalPages := (totalCount + limit - 1) / limit
	hasNext := page < totalPages
	hasPrev := page > 1

	// Prepare data for template
	data := map[string]interface{}{
		"User":         user,
		"Users":        users,
		"Stats":        stats,
		"SearchQuery":  searchQuery,
		"StatusFilter": statusFilter,
		"Page":         page,
		"TotalCount":   totalCount,
		"TotalPages":   totalPages,
		"HasNext":      hasNext,
		"HasPrev":      hasPrev,
	}

	// Render template
	templates.Temp.ExecuteTemplate(w, "admin_users", data)
}

// AdminBanUser handles banning/unbanning a user
func AdminBanUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Parse form
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get user ID and action from form
	userIdStr := r.FormValue("user_id")
	action := r.FormValue("action") // "ban" or "unban"

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Prevent admin from banning themselves
	if userId == user.UserId {
		http.Error(w, "Cannot ban yourself", http.StatusBadRequest)
		return
	}

	// Perform ban/unban action
	if action == "ban" {
		err = services.BanUserService(userId)
	} else if action == "unban" {
		err = services.UnbanUserService(userId)
	} else {
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Error performing action: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect back to users management
	http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
}

// AdminUpdateDiscussionStatus handles updating discussion status from admin panel
func AdminUpdateDiscussionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
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

	discussionId, err := strconv.Atoi(discussionIdStr)
	if err != nil {
		http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
		return
	}

	// Update discussion status (admin can update any discussion)
	err = services.UpdateDiscussionStatusService(discussionId, status, user.UserId, true)
	if err != nil {
		http.Error(w, "Error updating discussion status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get referrer to redirect back
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		referrer = "/admin/discussions"
	}

	http.Redirect(w, r, referrer, http.StatusSeeOther)
}

// AdminDeleteDiscussion handles deleting a discussion from admin panel
func AdminDeleteDiscussion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Extract ID from URL path or form
	var discussionId int
	if strings.HasSuffix(r.URL.Path, "/delete") {
		// Get from form
		discussionIdStr := r.FormValue("discussion_id")
		discussionId, err = strconv.Atoi(discussionIdStr)
		if err != nil {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}
	} else {
		// Get from URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}
		discussionId, err = strconv.Atoi(pathParts[len(pathParts)-1])
		if err != nil {
			http.Error(w, "Invalid discussion ID", http.StatusBadRequest)
			return
		}
	}

	// Delete discussion (admin can delete any discussion)
	err = services.DeleteDiscussionService(discussionId, user.UserId, true)
	if err != nil {
		http.Error(w, "Error deleting discussion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get referrer to redirect back
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		referrer = "/admin/discussions"
	}

	http.Redirect(w, r, referrer, http.StatusSeeOther)
}

// AdminDeleteMessage handles deleting a message from admin panel
func AdminDeleteMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if user is authenticated and is admin
	user, err := utils.GetUserFromRequest(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if user.Role != "admin" {
		http.Error(w, "Access denied. Admin privileges required.", http.StatusForbidden)
		return
	}

	// Extract ID from URL path or form
	var messageId int
	if strings.HasSuffix(r.URL.Path, "/delete") {
		// Get from form
		messageIdStr := r.FormValue("message_id")
		messageId, err = strconv.Atoi(messageIdStr)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
	} else {
		// Get from URL path
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
		messageId, err = strconv.Atoi(pathParts[len(pathParts)-1])
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}
	}

	// Delete message (admin can delete any message)
	err = services.DeleteMessageService(messageId, user.UserId, true)
	if err != nil {
		http.Error(w, "Error deleting message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get referrer to redirect back
	referrer := r.Header.Get("Referer")
	if referrer == "" {
		referrer = "/admin"
	}

	http.Redirect(w, r, referrer, http.StatusSeeOther)
}
