package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"server/middleware"
	"server/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func AddUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	//Changes the upvote number in posts/comments, inserts new addition to upvotes table
	var addUpvoteRequest models.UpvoteRequest
	err := json.NewDecoder(r.Body).Decode(&addUpvoteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(addUpvoteRequest)
	//Check username
	authUsername, _ := middleware.GetUsername(r)
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
		_, err = models.DataBase.Exec("UPDATE cvwo_assignment.posts SET upvotes = upvotes + 1 WHERE postid = $1", addUpvoteRequest.Postid)
		if err != nil {
			http.Error(w, "Failed to update upvotes in posts", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("INSERT INTO cvwo_assignment.upvotes (upvotes_postid, upvotes_username, datetime) VALUES ($1, $2, $3)", createdUpvote.Postid, createdUpvote.Username, createdUpvote.Datetime)
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
		_, err = models.DataBase.Exec("UPDATE cvwo_assignment.comments SET upvotes = upvotes + 1 WHERE commentid = $1", addUpvoteRequest.Commentid)
		if err != nil {
			http.Error(w, "Failed to update upvotes in comments", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("INSERT INTO cvwo_assignment.upvotes (upvotes_commentid, upvotes_username, datetime) VALUES ($1, $2, $3)", createdUpvote.Commentid, createdUpvote.Username, createdUpvote.Datetime)
		if err != nil {
			fmt.Println("DB error:", err)
			http.Error(w, "Failed to insert upvote into upvotes", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUpvote)
	fmt.Println("Upvotes inserted successfully")
}

func DeleteUpvoteHandler(w http.ResponseWriter, r *http.Request) {
	//Edits the posts/comments table upvotes and deletes the row from upvotes table
	var deleteUpvoteRequest models.UpvoteRequest
	err := json.NewDecoder(r.Body).Decode(&deleteUpvoteRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//Check username
	authUsername, _ := middleware.GetUsername(r)
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
		_, err = models.DataBase.Exec("UPDATE cvwo_assignment.posts SET upvotes = upvotes - 1 WHERE postid = $1", deleteUpvoteRequest.Postid)
		if err != nil {
			http.Error(w, "Failed to edit post", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("DELETE FROM cvwo_assignment.upvotes WHERE upvotes_postid = $1 AND upvotes_username = $2", deleteUpvoteRequest.Postid, deleteUpvoteRequest.Username)
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

		_, err = models.DataBase.Exec("UPDATE cvwo_assignment.comments SET upvotes = upvotes - 1 WHERE commentid = $1", deleteUpvoteRequest.Commentid)
		if err != nil {
			http.Error(w, "Failed to edit comment", http.StatusInternalServerError)
			return
		}
		_, err = models.DataBase.Exec("DELETE FROM cvwo_assignment.upvotes WHERE upvotes_commentid = $1 AND upvotes_username = $2", deleteUpvoteRequest.Commentid, deleteUpvoteRequest.Username)
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
	username, _ := middleware.GetUsername(r)
	row := models.DataBase.QueryRow("SELECT upvotes_postid, upvotes_username, datetime FROM cvwo_assignment.upvotes WHERE upvotes_postid = $1 AND upvotes_username = $2", postID, username)
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
		fmt.Println("No upvote:", err)
		http.Error(w, "Upvote does not exist", http.StatusInternalServerError)
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
	username, _ := middleware.GetUsername(r)
	row := models.DataBase.QueryRow("SELECT upvotes_commentid, upvotes_username, datetime FROM cvwo_assignment.upvotes WHERE upvotes_commentid = $1 AND upvotes_username = $2", commentID, username)
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
		fmt.Println("No upvote:", err)
		http.Error(w, "Upvote does not exist", http.StatusInternalServerError)
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
		query = "SELECT EXISTS(SELECT 1 FROM cvwo_assignment.upvotes WHERE upvotes_postid = $1 AND upvotes_username = $2 LIMIT 1)"
	} else {
		query = "SELECT EXISTS(SELECT 1 FROM cvwo_assignment.upvotes WHERE upvotes_commentid = $1 AND upvotes_username = $2 LIMIT 1)"
	}
	err := models.DataBase.QueryRow(query, ID, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func dateTimeConverter(datetimeRaw []uint8) time.Time {
	datetimeStr := string(datetimeRaw)
	datetime, err := time.Parse("2006-01-02T15:04:05.9999999Z", datetimeStr)
	if err != nil {
		fmt.Println("Error parsing datetime:", err)
		return time.Now()
	}
	return datetime
}
