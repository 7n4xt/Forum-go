package routes

import (
	"forum-go/controllers"
	"net/http"
)

func MessageRouter(mux *http.ServeMux) {
	mux.HandleFunc("/messages", controllers.GetMessages)
	mux.HandleFunc("/messages/", controllers.GetMessageByID)
	mux.HandleFunc("/messages", controllers.CreateMessage)
	mux.HandleFunc("/messages/", controllers.UpdateMessage)
	mux.HandleFunc("/messages/", controllers.DeleteMessage)
}
