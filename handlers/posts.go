package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"server/models"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// ************Posts handlers*****************
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	//Extract sorting options from query parameters +
	sortBy := r.URL.Query().Get("sort_by")
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("perPage")

	query := "SELECT * FROM cvwo_assignment.posts"
	var args []interface{}

	//Sorting options
	addedSortingQuery := addSorting(query, sortBy)
	// Pagination
	addedPaginationQuery, args := addPagination(addedSortingQuery, pageStr, perPageStr, 0)

	//Query DB
	rows, err := models.DataBase.Query(addedPaginationQuery, args...)
	if err != nil {
		fmt.Println("sql err:", err)
		fmt.Println(addedPaginationQuery)
		fmt.Println(args...)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := scanPosts(rows)
	if err != nil {
		fmt.Println("Error scanning:", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
	fmt.Println("getPosts handler worked")
}

func GetPostByIdHandler(w http.ResponseWriter, r *http.Request) { //Get post by post id
	fmt.Println("Entered GetPost by id handler")
	params := mux.Vars(r)
	postIdstr := params["postid"]
	postID, err := strconv.Atoi(postIdstr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	fmt.Println(postID)
	post, err := getPostById(postID)
	if err != nil {
		fmt.Println("Err getting post by id:", err)
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
	fmt.Println("getPost handler worked")
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	//Recieves post request then inserts it into post table and insert tags into posttags table
	var newPost models.PostRequest
	err := json.NewDecoder(r.Body).Decode(&newPost)
	if err != nil {
		http.Error(w, "Invalid body request", http.StatusBadRequest)
		return
	}

	createdPost := models.InitPost(newPost.Username, newPost.Tags, newPost.Title, newPost.Description)

	tagsJSON, err := json.Marshal(newPost.Tags)
	if err != nil {
		http.Error(w, "Failed to encode tags", http.StatusInternalServerError)
		return
	}
	var postID int64
	err = models.DataBase.QueryRow("INSERT INTO cvwo_assignment.posts (username, tag, title, description, upvotes, datetime) VALUES ($1, $2, $3, $4, $5, $6) RETURNING postid", createdPost.Username, string(tagsJSON), createdPost.Title, createdPost.Description, createdPost.Upvote, createdPost.Datetime).Scan(&postID)
	if err != nil {
		http.Error(w, "Failed to create post or get post ID", http.StatusInternalServerError)
		return
	}

	createdPost.Postid = postID

	err = insertTagsForPost(postID, newPost.Tags)
	if err != nil {
		http.Error(w, "Failed to insert tags to create psot", http.StatusBadRequest)
		return
	}
	fmt.Println("Tags inserted successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdPost)
	fmt.Println("createPosts handler worked")
}

func EditPostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postIDStr := params["postid"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var updatedPost models.PostRequest
	err = json.NewDecoder(r.Body).Decode(&updatedPost)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tagsJSON, err := json.Marshal(updatedPost.Tags)
	if err != nil {
		http.Error(w, "Failed to encode tags", http.StatusInternalServerError)
		return
	}

	_, err = models.DataBase.Exec("UPDATE cvwo_assignment.posts SET tag = $1, title = $2, description = $3 WHERE postid = $4", tagsJSON, updatedPost.Title, updatedPost.Description, postID)
	if err != nil {
		http.Error(w, "Failed to edit post", http.StatusInternalServerError)
		return
	}
	editedPost := models.EditPost(int64(postID), updatedPost.Username, updatedPost.Tags, updatedPost.Title, updatedPost.Description)

	err = editTagsForPost(int(editedPost.Postid), editedPost.Tags)
	if err != nil {
		http.Error(w, "Failed to edit tags", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editedPost)
	fmt.Println("editPost handler worked")
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postIDStr := params["postid"]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	_, err = models.DataBase.Exec("DELETE FROM cvwo_assignment.posts WHERE postid = $1", postID)
	if err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	err = deleteTagsForPost(postID)
	if err != nil {
		http.Error(w, "Failed to delete tags", http.StatusInternalServerError)
		return
	}
	fmt.Println("Deleted post")
	w.WriteHeader(http.StatusNoContent)
}

//***************Abstracted Functions***************

func scanPosts(rows *sql.Rows) ([]models.Post, error) {
	//scans rows of posts into slice of posts
	var posts []models.Post
	for rows.Next() {
		var post models.Post
		//Tags and datetime needs additional processing
		var tagsJSON []byte
		var datetimeRaw []uint8

		err := rows.Scan(
			&post.Postid,
			&post.Username,
			&tagsJSON,
			&post.Title,
			&post.Description,
			&post.Upvote,
			&datetimeRaw,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		var tags []string
		if err := json.Unmarshal(tagsJSON, &tags); err != nil {
			fmt.Println("Error unmarshaling JSON:", err)
			return nil, err
		}
		datetime := dateTimeConverter(datetimeRaw)
		post.Tags = tags
		post.Datetime = datetime
		posts = append(posts, post)
	}
	return posts, nil
}

func getPostById(postID int) (models.Post, error) {
	fmt.Println(postID)
	row := models.DataBase.QueryRow("SELECT * FROM cvwo_assignment.posts WHERE postid = $1", postID)

	var post models.Post
	//Tags and datetime needs additional processing
	var tagsJSON []byte
	var datetimeRaw []uint8

	err := row.Scan(&post.Postid, &post.Username, &tagsJSON, &post.Title, &post.Description, &post.Upvote, &datetimeRaw)
	if err == sql.ErrNoRows {
		log.Println("Post not found")
		return models.Post{}, err
	} else if err != nil {
		log.Printf("Error scanning post: %v", err)
		return models.Post{}, err
	}
	//Additional processing of tags and datetime
	var tags []string
	if err := json.Unmarshal(tagsJSON, &tags); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return models.Post{}, err
	}
	datetime := dateTimeConverter(datetimeRaw)
	post.Tags = tags
	post.Datetime = datetime
	return post, nil
}

func addSorting(sqlquery string, sortBy string) string {
	switch sortBy {
	case "time":
		sqlquery += " ORDER BY datetime DESC"
	case "upvotes":
		sqlquery += " ORDER BY upvotes DESC"
	default:
		sqlquery += " ORDER BY datetime DESC"
	}
	return sqlquery
}

func addPagination(query string, pageStr string, perPageStr string, len int) (string, []interface{}) {
	var args []interface{}
	if pageStr != "" && perPageStr != "" {
		pageNumber, _ := strconv.Atoi(pageStr)
		itemsPerPage, _ := strconv.Atoi(perPageStr)

		if pageNumber <= 0 {
			pageNumber = 1
		}

		if itemsPerPage <= 0 {
			itemsPerPage = 5
		}
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len+1, len+2)
		args = append(args, itemsPerPage, (pageNumber-1)*itemsPerPage)

		return query, args
	}
	return query, nil
}
