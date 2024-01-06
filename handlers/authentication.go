package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"
	"server/sessions"
	"time"

	_ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	s "github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// *************** Authentication Handlers ********************
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON input
	var loginRequest models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate required fields
	if loginRequest.Username == "" || loginRequest.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	//Query database for user
	var hashedPassword string
	err = models.DataBase.QueryRow("SELECT password FROM cvwo_assignment.users WHERE username = $1", loginRequest.Username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username", http.StatusUnauthorized)
			return
		}
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Compare the provided password with the hashed password from the database
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	fmt.Println("Login Success.... Creating token")
	//Authentication successful, create jwt token
	token, err := sessions.CreateJWT(loginRequest.Username)
	if err != nil {
		http.Error(w, "Error creating jwt token", http.StatusInternalServerError)
		return
	}
	//Creating session to store token
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	session.Options = &s.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   "https://web-forum-jmof.onrender.com/",
	}
	session.Values["jwt-token"] = token
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	//Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login successful"}`))
	fmt.Println("loginHandler worked")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered logout handler")
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	// Clear jwt
	delete(session.Values, "jwt-token")
	// Set MaxAge to -1 for immediate expiration
	session.Options.MaxAge = -1

	// Save the session with updated MaxAge
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logout successful"}`))
	fmt.Println("logoutHandler worked")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//Parse JSON input
	var registerRequest models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Validate requried fields
	if registerRequest.Username == "" || registerRequest.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	var existingUsername string
	err = models.DataBase.QueryRow("SELECT username FROM cvwo_assignment.users WHERE username = $1", registerRequest.Username).Scan(&existingUsername)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error at query row:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if existingUsername != "" {
		http.Error(w, "Username is already taken", http.StatusConflict)
		return
	}
	//Hash password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//Insert user into database
	_, err = models.DataBase.Exec("INSERT INTO cvwo_assignment.users (username, password, created_at) VALUES ($1, $2, $3)",
		registerRequest.Username, hashPassword, time.Now())
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//Response headers
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Registration successful"}`))
	fmt.Println("registerHandler worked")
}
