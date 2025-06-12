package routes

import (
	"forum-go/controllers"
	"net/http"
)

// AuthRouter registers all authentication related routes to the provided mux
func AuthRouter(mux *http.ServeMux) {
	mux.HandleFunc("/signup", controllers.PageSignUp)
	mux.HandleFunc("/signup/form", controllers.FormSignup)
	mux.HandleFunc("/login", controllers.PageLogin)
	mux.HandleFunc("/login/form", controllers.FormLogin)
	mux.HandleFunc("/logout", controllers.Logout)
}
