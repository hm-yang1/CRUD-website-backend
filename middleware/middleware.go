package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"server/models"
	"server/sessions"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// ******************** Authorisation checks ***********************
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	// Checks if session present
	fmt.Println("Session Checker fired")
	session, _ := sessions.Store.Get(r, "session")
	token, ok := session.Values["jwt-token"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	_, err := sessions.ParseJWT(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Session present"}`))
}

func AuthRequired(handler http.HandlerFunc) http.HandlerFunc {
	// middleware to check if user has a session. If present, allows handler to go through
	fmt.Println("AuthRequired fired")
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		token, ok := session.Values["jwt-token"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		_, err := sessions.ParseJWT(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func IsCommentEditAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	//Checks if comment edit is allowed by comparing session username and comment username
	return func(w http.ResponseWriter, r *http.Request) {
		//Get username from jwt in session
		authenticatedUsername, err := GetUsername(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		//Parse comment id from request parameters
		params := mux.Vars(r)
		commentIDStr := params["commentid"]
		commentID, err := strconv.Atoi(commentIDStr)
		if err != nil {
			http.Error(w, "Invalid comment ID", http.StatusBadRequest)
			return
		}
		var commentUsername string
		err = models.DataBase.QueryRow("SELECT comments_username FROM cvwo_assignment.comments WHERE commentid = $1", commentID).Scan(&commentUsername)
		if err != nil {
			http.Error(w, "Unable to match comment to username", http.StatusInternalServerError)
			return
		}
		if authenticatedUsername != commentUsername {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func IsPostEditAuthorized(handler http.HandlerFunc) http.HandlerFunc {
	//Checks if post edit is allowed by comparing session username and post username
	return func(w http.ResponseWriter, r *http.Request) {
		//Get username from session
		authenticatedUsername, err := GetUsername(r)
		if err != nil {
			http.Error(w, "No username", http.StatusUnauthorized)
			return
		}
		fmt.Println(authenticatedUsername)
		//Parse post id from request parameters
		params := mux.Vars(r)
		postIDStr := params["postid"]
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		var postUsername string
		err = models.DataBase.QueryRow("SELECT username FROM cvwo_assignment.posts WHERE postid = $1", postID).Scan(&postUsername)
		if err != nil {
			http.Error(w, "Unable to find post", http.StatusInternalServerError)
			return
		}
		if authenticatedUsername != postUsername {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
		handler.ServeHTTP(w, r)
	}
}

func GetUsername(r *http.Request) (string, error) {
	//Gets username from jwt in session
	session, err := sessions.Store.Get(r, "session")
	token, ok := session.Values["jwt-token"].(string)
	if !ok {
		err = errors.New("jwt tokent from session")
		return "", err
	}
	claims, err := sessions.ParseJWT(token)
	if err != nil {
		return "", err
	}
	authenticatedUsername, ok := claims["username"].(string)
	if !ok {
		err = errors.New("Failed to get username from jwt")
		return "", err
	}
	return authenticatedUsername, nil
}
