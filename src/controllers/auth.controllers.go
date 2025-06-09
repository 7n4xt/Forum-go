package controllers

import (
	"fmt"
	"forum-go/services"
	"forum-go/templates"
	"net/http"
)

func PageSignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page signup called")
	templates.Temp.ExecuteTemplate(w, "signup", nil)
}

func FormSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Form signup called")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, IsCreateErr := services.SignupService(username, password, email)
	if IsCreateErr != nil {
		fmt.Println("Error during signup:", IsCreateErr)
		http.Error(w, IsCreateErr.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func PageLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Page login called")
	templates.Temp.ExecuteTemplate(w, "login", nil)
}
func FormLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Form login called")
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	password := r.FormValue("password")
	username := r.FormValue("username")

	token, err := services.LoginService(password, username)
	if err != nil {
		fmt.Println("Error during login:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Login successful, token:", token)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    token,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/discussions", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Delete the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
