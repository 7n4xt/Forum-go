package routes

import (
	"forum-go/controllers"
	"forum-go/utils"
	"net/http"
)

// DiscussionRouter registers all discussion related routes
func DiscussionRouter(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("/discussions", controllers.GetAllDiscussions)
	mux.HandleFunc("/discussions/", controllers.GetDiscussionByID)

	// Protected routes (require authentication)
	mux.HandleFunc("/discussions/new", utils.RequireAuthentication(controllers.NewDiscussionPage))
	mux.HandleFunc("/discussions/create", utils.RequireAuthentication(controllers.CreateDiscussionForm))
	mux.HandleFunc("/discussions/update", utils.RequireAuthentication(controllers.UpdateDiscussionContent))
	mux.HandleFunc("/discussions/status", utils.RequireAuthentication(controllers.UpdateDiscussionStatus))
	mux.HandleFunc("/discussions/delete", utils.RequireAuthentication(controllers.DeleteDiscussion))
	mux.HandleFunc("/discussions/delete/", utils.RequireAuthentication(controllers.DeleteDiscussion))
}
