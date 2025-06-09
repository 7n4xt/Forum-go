package routes

import (
	"forum-go/controllers"
	"forum-go/utils"
	"net/http"
)

// MessageRouter registers all message related routes
func MessageRouter(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("/messages", controllers.GetMessages)
	mux.HandleFunc("/messages/view/", controllers.GetMessageByID)

	// Protected routes (require authentication)
	mux.HandleFunc("/messages/create", utils.RequireAuthentication(controllers.CreateMessage))
	mux.HandleFunc("/messages/create/", utils.RequireAuthentication(controllers.CreateMessage))
	mux.HandleFunc("/messages/update/", utils.RequireAuthentication(controllers.UpdateMessage))
	mux.HandleFunc("/messages/delete/", utils.RequireAuthentication(controllers.DeleteMessage))
}
