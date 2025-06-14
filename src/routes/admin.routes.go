package routes

import (
	"forum-go/controllers"
	"forum-go/utils"
	"net/http"
)

// AdminRouter registers all admin related routes
func AdminRouter(mux *http.ServeMux) {
	// All admin routes require authentication and admin privileges
	mux.HandleFunc("/admin", utils.RequireAuthentication(controllers.AdminDashboard))
	mux.HandleFunc("/admin/", utils.RequireAuthentication(controllers.AdminDashboard))

	// User management routes
	mux.HandleFunc("/admin/users", utils.RequireAuthentication(controllers.AdminManageUsers))
	mux.HandleFunc("/admin/users/ban", utils.RequireAuthentication(controllers.AdminBanUser))

	// Discussion management routes
	mux.HandleFunc("/admin/discussions", utils.RequireAuthentication(controllers.AdminManageDiscussions))
	mux.HandleFunc("/admin/discussions/status", utils.RequireAuthentication(controllers.AdminUpdateDiscussionStatus))
	mux.HandleFunc("/admin/discussions/delete", utils.RequireAuthentication(controllers.AdminDeleteDiscussion))
	mux.HandleFunc("/admin/discussions/delete/", utils.RequireAuthentication(controllers.AdminDeleteDiscussion))

	// Message management routes
	mux.HandleFunc("/admin/messages/delete", utils.RequireAuthentication(controllers.AdminDeleteMessage))
	mux.HandleFunc("/admin/messages/delete/", utils.RequireAuthentication(controllers.AdminDeleteMessage))
}
