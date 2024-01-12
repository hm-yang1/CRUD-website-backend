package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/models"

	_ "github.com/lib/pq"
)

// **************Search handlers*********************
// Currently only got search posts
func GetSearchedPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered GetSearchedPosts")
	//Get queries
	searchQuery := r.URL.Query().Get("query")
	sortBy := r.URL.Query().Get("sort_by")
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("perPage")
	fmt.Println(searchQuery)
	query := `
		SELECT * 
		FROM cvwo_assignment.posts
		WHERE title LIKE $1 OR description LIKE $2
	`
	//Sorting options
	addedSortingQuery := addSorting(query, sortBy)
	//Add pagination
	addedPaginationQuery, args := addPagination(addedSortingQuery, pageStr, perPageStr, 2)
	fullArgs := append([]interface{}{"%" + searchQuery + "%", "%" + searchQuery + "%"}, args...)
	rows, err := models.DataBase.Query(addedPaginationQuery, fullArgs...)
	if err != nil {
		fmt.Println("Error in searchHandler:", err)
		http.Error(w, "Error searching for posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	posts, err := scanPosts(rows)
	if err != nil {
		fmt.Println("Error scanning posts:", err)
		http.Error(w, "Error scanning posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
