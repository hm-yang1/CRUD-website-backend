package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//************** Comments Handlers *************

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	//Get postID
	params := mux.Vars(r)
	postIdstr := params["postid"]
	postID, err := strconv.Atoi(postIdstr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	//Sorting
	sortBy := r.URL.Query().Get("sort_by")
	query := "SELECT * FROM cvwo_assignment.comments WHERE comments_postid = $1"
	// var args []interface{}
	//Added sorting options
	addedSortingQuery := addSorting(query, sortBy)
	//Get all comments with same postid
	rows, err := models.DataBase.Query(addedSortingQuery, postID)
	if err != nil {
		log.Printf("Error getting comment: %v", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	//Get slice of comments
	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		var datetimeRaw []uint8

		err := rows.Scan(
			&comment.Commentid,
			&comment.Postid,
			&comment.Username,
			&comment.Description,
			&comment.Upvote,
			&datetimeRaw,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		datetimeStr := string(datetimeRaw)
		fmt.Println(datetimeStr)
		datetime, err := time.Parse("2006-01-02T15:04:05.9999999Z", datetimeStr)
		if err != nil {
			// Handle the error, e.g., log it or return an error
			fmt.Println("Error parsing datetime:", err)
			return
		}
		comment.Datetime = datetime

		comments = append(comments, comment)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
	fmt.Println("getComment handler worked")
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered GetComment handler")
	params := mux.Vars(r)
	commentIdStr := params["commentid"]
	commentID, err := strconv.Atoi(commentIdStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var getComment models.Comment
	var datetimeRaw []uint8

	row := models.DataBase.QueryRow("SELECT * FROM cvwo_assignment.comments WHERE commentid = $1", commentID)
	err = row.Scan(
		&getComment.Commentid,
		&getComment.Postid,
		&getComment.Username,
		&getComment.Description,
		&getComment.Upvote,
		&datetimeRaw,
	)
	if err != nil {
		log.Printf("Error scanning: %v", err)
		http.Error(w, "Failed to scan edited comment", http.StatusInternalServerError)
		return
	}
	datetimeStr := string(datetimeRaw)
	datetime, err := time.Parse("2006-01-02T15:04:05.9999999Z", datetimeStr)
	if err != nil {
		// Handle the error, e.g., log it or return an error
		fmt.Println("Error parsing datetime:", err)
		return
	}
	getComment.Datetime = datetime

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getComment)
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered createComment handler")
	//Decode comment
	var newComment models.CommentRequest
	err := json.NewDecoder(r.Body).Decode(&newComment)
	if err != nil {
		http.Error(w, "Invalid body request", http.StatusBadRequest)
		return
	}
	//Insert comment into DB
	var createdComment models.Comment
	createdComment = models.InitComment(newComment.Postid, newComment.Username, newComment.Description)
	err = models.DataBase.QueryRow("INSERT INTO cvwo_assignment.comments (comments_postid, comments_username, description, upvotes, datetime) VALUES ($1, $2, $3, $4, $5) RETURNING commentid", createdComment.Postid, createdComment.Username, createdComment.Description, createdComment.Upvote, createdComment.Datetime).Scan(&createdComment.Commentid)
	if err != nil {
		log.Printf("Error inserting: %v", err)
		http.Error(w, "Failed to create comment or get comment ID", http.StatusInternalServerError)
		return
	}
	//Return new comment
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdComment)
	fmt.Println("createComment worked")
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Parse commentIDStr to int
	fmt.Println("Entered EditComment")
	params := mux.Vars(r)
	commentIDStr := params["commentid"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var updatedComment models.CommentRequest
	err = json.NewDecoder(r.Body).Decode(&updatedComment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Execute the UPDATE query to edit the existing comment
	_, err = models.DataBase.Exec("UPDATE cvwo_assignment.comments SET description = $1 WHERE commentid = $2", updatedComment.Description, commentID)
	if err != nil {
		http.Error(w, "Failed to edit comment", http.StatusInternalServerError)
		return
	}
	// Query and return edited comment
	var editedComment models.Comment
	var datetimeRaw []uint8

	row := models.DataBase.QueryRow("SELECT * FROM cvwo_assignment.comments WHERE commentid = $1", commentID)
	err = row.Scan(
		&editedComment.Commentid,
		&editedComment.Postid,
		&editedComment.Username,
		&editedComment.Description,
		&editedComment.Upvote,
		&datetimeRaw,
	)
	if err != nil {
		log.Printf("Error scanning: %v", err)
		http.Error(w, "Failed to scan edited comment", http.StatusInternalServerError)
		return
	}
	datetimeStr := string(datetimeRaw)
	datetime, err := time.Parse("2006-01-02T15:04:05.9999999Z", datetimeStr)
	if err != nil {
		// Handle the error, e.g., log it or return an error
		fmt.Println("Error parsing datetime:", err)
		return
	}
	editedComment.Datetime = datetime

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editedComment)
	fmt.Println("editCommentHandler worked")
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	//Extract comment id
	params := mux.Vars(r)
	commentIDStr := params["commentid"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	//Execute DELETE query in DB
	_, err = models.DataBase.Exec("DELETE FROM cvwo_assignment.comments WHERE commentid = $1", commentID)
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
