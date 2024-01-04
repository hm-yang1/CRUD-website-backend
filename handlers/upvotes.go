package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server/models"
	"server/sessions"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func AddUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	var addUpvoteRequest models.UpvoteRequest
	err := json.NewDecoder(r.Body).Decode(&addUpvoteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(addUpvoteRequest)
	//Check username
	authUsername := checkUsername(r)
	if authUsername != addUpvoteRequest.Username {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var createdUpvote models.Upvote
	createdUpvote = models.Upvote{
		Postid:    addUpvoteRequest.Postid,
		Commentid: addUpvoteRequest.Commentid,
		Username:  addUpvoteRequest.Username,
		Datetime:  time.Now(),
	}
	if addUpvoteRequest.Commentid == 0 && addUpvoteRequest.Postid > 0 {
		//Check if upvote already exist
		existsPost, err := rowExists(int(addUpvoteRequest.Postid), addUpvoteRequest.Username, true)
		if err != nil {
			fmt.Println("upvotes table query failed", err)
			http.Error(w, "upvotes query failed", http.StatusInternalServerError)
			return
		}
		if existsPost {
			http.Error(w, "Already Upvoted", http.StatusUnauthorized)
			return
		} else {
			fmt.Println("Good to go")
		}
		_, err = models.DataBase.Exec("UPDATE posts SET upvotes = upvotes + 1 WHERE postid = ?", addUpvoteRequest.Postid)
		if err != nil {
			http.Error(w, "Failed to update upvotes in posts", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("INSERT INTO upvotes (upvotes_postid, upvotes_username, datetime) VALUES (?, ?, ?)", createdUpvote.Postid, createdUpvote.Username, createdUpvote.Datetime)
		if err != nil {
			fmt.Println("DB error:", err)
			http.Error(w, "Failed to insert upvote into upvotes", http.StatusInternalServerError)
			return
		}
	} else {
		existsPost, err := rowExists(int(addUpvoteRequest.Commentid), addUpvoteRequest.Username, false)
		if err != nil {
			fmt.Println("upvotes table query failed", err)
			http.Error(w, "upvotes query failed", http.StatusInternalServerError)
			return
		}
		if existsPost {
			http.Error(w, "Already Upvoted", http.StatusUnauthorized)
			return
		} else {
			fmt.Println("Good to go")
		}
		_, err = models.DataBase.Exec("UPDATE comments SET upvotes = upvotes + 1 WHERE commentid = ?", addUpvoteRequest.Commentid)
		if err != nil {
			http.Error(w, "Failed to update upvotes in comments", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("INSERT INTO upvotes (upvotes_commentid, upvotes_username, datetime) VALUES (?, ?, ?)", createdUpvote.Commentid, createdUpvote.Username, createdUpvote.Datetime)
		if err != nil {
			fmt.Println("DB error:", err)
			http.Error(w, "Failed to insert upvote into upvotes", http.StatusInternalServerError)
			return
		}
	}
	//Add upvote to upvotes table
	fmt.Println(createdUpvote)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUpvote)
	fmt.Println("Upvotes inserted successfully")
}

func DeleteUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	//Edits the posts/comments table upvotes by minus 1 and deletes the row from upvotes table
	//Bends rest principles but I want thoose 2 things to happen in 1 function
	var deleteUpvoteRequest models.UpvoteRequest
	err := json.NewDecoder(r.Body).Decode(&deleteUpvoteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Check username
	authUsername := checkUsername(r)
	if authUsername != deleteUpvoteRequest.Username {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if deleteUpvoteRequest.Commentid == 0 && deleteUpvoteRequest.Postid > 0 {
		existsPost, err := rowExists(int(deleteUpvoteRequest.Postid), deleteUpvoteRequest.Username, true)
		if err != nil {
			fmt.Println("upvotes table query failed", err)
			http.Error(w, "upvotes query failed", http.StatusInternalServerError)
			return
		}
		if existsPost {
			fmt.Println("Good to go")
		} else {
			http.Error(w, "No upvote", http.StatusUnauthorized)
			return
		}
		_, err = models.DataBase.Exec("UPDATE posts SET upvotes = upvotes - 1 WHERE postid = ?", deleteUpvoteRequest.Postid)
		if err != nil {
			http.Error(w, "Failed to edit post", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("DELETE FROM upvotes WHERE upvotes_postid = ? AND upvotes_username = ?", deleteUpvoteRequest.Postid, deleteUpvoteRequest.Username)
		if err != nil {
			fmt.Println("Delete Error:", err)
			http.Error(w, "Failed to delete upvote", http.StatusInternalServerError)
			return
		}
	} else {
		existsPost, err := rowExists(int(deleteUpvoteRequest.Commentid), deleteUpvoteRequest.Username, false)
		if err != nil {
			fmt.Println("upvotes table query failed", err)
			http.Error(w, "upvotes query failed", http.StatusInternalServerError)
			return
		}
		if existsPost {
			fmt.Println("Good to go")
		} else {
			http.Error(w, "No upvote", http.StatusUnauthorized)
			return
		}

		_, err = models.DataBase.Exec("UPDATE comments SET upvotes = upvotes - 1 WHERE commentid = ?", deleteUpvoteRequest.Commentid)
		if err != nil {
			http.Error(w, "Failed to edit comment", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("DELETE FROM upvotes WHERE upvotes_commentid = ? AND upvotes_username = ?", deleteUpvoteRequest.Commentid, deleteUpvoteRequest.Username)
		if err != nil {
			fmt.Println("Delete Error:", err)
			http.Error(w, "Failed to delete upvote", http.StatusInternalServerError)
			return
		}
	}
	fmt.Println("Deleted upvote")
	w.WriteHeader(http.StatusNoContent)
}

func GetPostUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postIDStr := params["postid"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	username := checkUsername(r)
	row := models.DataBase.QueryRow("SELECT upvotes_postid, upvotes_username, datetime FROM upvotes WHERE upvotes_postid = ? AND upvotes_username = ?", postID, username)
	var upvote models.Upvote
	upvote = models.Upvote{
		Postid:    -1,
		Commentid: -1,
		Username:  "",
		Datetime:  time.Now(),
	}
	var datetimeRaw []uint8
	err = row.Scan(&upvote.Postid, &upvote.Username, &datetimeRaw)
	if err == sql.ErrNoRows {
		w.Header().Set("Content=Type", "application/json")
		json.NewEncoder(w).Encode(upvote)
		return
	} else if err != nil {
		fmt.Println("Error scanning upvote", err)
		http.Error(w, "Failed to fetch upvote", http.StatusInternalServerError)
		return
	}
	upvote.Datetime = dateTimeConverter(datetimeRaw)
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(upvote)
}

func GetCommentUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	commentIDStr := params["commentid"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	username := checkUsername(r)
	row := models.DataBase.QueryRow("SELECT upvotes_commentid, upvotes_username, datetime FROM upvotes WHERE upvotes_commentid = ? AND upvotes_username = ?", commentID, username)
	var upvote models.Upvote
	upvote = models.Upvote{
		Postid:    -1,
		Commentid: -1,
		Username:  "",
		Datetime:  time.Now(),
	}
	var datetimeRaw []uint8
	err = row.Scan(&upvote.Commentid, &upvote.Username, &datetimeRaw)
	if err == sql.ErrNoRows {
		w.Header().Set("Content=Type", "application/json")
		json.NewEncoder(w).Encode(upvote)
		return
	} else if err != nil {
		fmt.Println("Error scanning upvote", err)
		http.Error(w, "Failed to fetch upvote", http.StatusInternalServerError)
		return
	}
	upvote.Datetime = dateTimeConverter(datetimeRaw)
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(upvote)
}

// ************Helper Functions*************
func rowExists(ID int, username string, post bool) (bool, error) {
	var exists bool
	var query string
	if post {
		query = "SELECT EXISTS(SELECT 1 FROM upvotes WHERE upvotes_postid = ? AND upvotes_username = ? LIMIT 1)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM upvotes WHERE upvotes_commentid = ? AND upvotes_username = ? LIMIT 1)"
	}
	err := models.DataBase.QueryRow(query, ID, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func checkUsername(r *http.Request) string {
	session, _ := sessions.Store.Get(r, "session")
	token, ok := session.Values["jwt-token"].(string)
	if !ok {
		return ""
	}
	claims, err := sessions.ParseJWT(token)
	if err != nil {
		return ""
	}
	authenticatedUsername, ok := claims["username"].(string)
	if !ok {
		return ""
	}
	return authenticatedUsername
}

func dateTimeConverter(datetimeRaw []uint8) time.Time {
	datetimeStr := string(datetimeRaw)
	datetime, err := time.Parse("2006-01-02 15:04:05", datetimeStr)
	if err != nil {
		// Handle the error, e.g., log it or return an error
		fmt.Println("Error parsing datetime:", err)
		return time.Now() //
	}
	return datetime
}
