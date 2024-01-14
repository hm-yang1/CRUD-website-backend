package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"server/models"
	"server/sessions"
	"time"

	s "github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// *************** Authentication Handlers ********************
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parses login request, compares passwords then creates session upon success
	var loginRequest models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if loginRequest.Username == "" || loginRequest.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
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
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	frontendUrl := os.Getenv("FRONTEND_URL")
	if frontendUrl == "" {
		log.Fatal("allowed origins empty")
	}
	session.Options = &s.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Domain:   frontendUrl,
	}
	session.Values["jwt-token"] = token
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Login successful"}`))
	fmt.Println("loginHandler worked")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	//Checks for session then clears the session
	session, err := sessions.Store.Get(r, "session")
	if err != nil {
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}
	// Clear jwt
	delete(session.Values, "jwt-token")
	// Set MaxAge to -1 for immediate expiration
	session.Options = &s.Options{
		Path:     "/",
		MaxAge:   -1, // Immediate expiration for logout
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true, // Set to true if your site is served over HTTPS
	}

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
	//Parse register request, checks requried fields, insert user into db
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
