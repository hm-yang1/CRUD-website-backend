package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/models"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	// _ "github.com/go-sql-driver/mysql"
)

// **********Tags handlers*****************
func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM cvwo_assignment.tags"
	//Query DB
	rows, err := models.DataBase.Query(query)
	if err != nil {
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tags []models.Tag

	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.Tagid,
			&tag.Name,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		tags = append(tags, tag)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
	fmt.Println("getTags handler worked")
}

func GetFilteredPostsHandler(w http.ResponseWriter, r *http.Request) {
	//Get posts filtered by tags
	//Get queries
	tagNames := r.URL.Query()["tag"]
	sortBy := r.URL.Query().Get("sort_by")
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("perPage")
	fmt.Println("tagNames:", tagNames)
	filteredPosts, err := getFilteredPosts(tagNames, sortBy, pageStr, perPageStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error retrieving filtered posts", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(filteredPosts)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// ******Functions used to manipulate the posttags table**********
func insertTagsForPost(postID int64, tagNames []string) error {
	// Retrieve tag ids
	tagIDs := make([]int, 0, len(tagNames))
	for _, tagName := range tagNames {
		var tagID int
		err := models.DataBase.QueryRow("SELECT tagsid FROM cvwo_assignment.tags WHERE name = $1", tagName).Scan(&tagID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		tagIDs = append(tagIDs, tagID)
	}

	fmt.Println("Selected tags", tagIDs)

	// Insert entries into posttags Table
	for _, tagID := range tagIDs {
		_, err := models.DataBase.Exec("INSERT INTO cvwo_assignment.posttags (posttags_tagid, posttags_postid) VALUES ($1, $2)", tagID, postID)
		if err != nil {
			return err
		}
	}
	fmt.Println("Inserted into posttags")
	return nil
}

func editTagsForPost(postID int, newTagNames []string) error {
	// Delete original tags
	err := deleteTagsForPost(postID)
	if err != nil {
		fmt.Println("Failed to delete", err)
		return err
	}
	// Retrieve and insert new tag IDs
	for _, tagName := range newTagNames {
		var tagID int
		err := models.DataBase.QueryRow("SELECT tagsid FROM cvwo_assignment.tags WHERE name = $1", tagName).Scan(&tagID)
		if err != nil {
			fmt.Println("Failed to edit", err)
			return err
		}
		_, err = models.DataBase.Exec("INSERT INTO cvwo_assignment.posttags (posttags_tagid, posttags_postid) VALUES ($1, $2)", tagID, postID)
		if err != nil {
			fmt.Println("Failed to edit", err)
			return err
		}
	}
	return nil
}

func deleteTagsForPost(postID int) error {
	tx, err := models.DataBase.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM cvwo_assignment.posttags WHERE posttags_postid = $1", postID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// ***************Functions used by the tag filtering handler**************
func getTagIdFromTagName(tagNames []string) ([]int, error) {
	tagIDs := make([]int, 0, len(tagNames))
	for _, tagName := range tagNames {
		var tagID int
		err := models.DataBase.QueryRow("SELECT tagsid FROM cvwo_assignment.tags WHERE name = $1", tagName).Scan(&tagID)
		if err != nil {
			fmt.Println(err)
			return []int{}, err
		}
		tagIDs = append(tagIDs, tagID)
	}
	return tagIDs, nil
}

func getPostIdsFromTagId(tagIds []int) ([]int, error) {
	if len(tagIds) == 0 {
		return nil, nil // No tag IDs provided, return an empty result
	}
	// Construct placeholders for the IN clause
	placeholders := make([]string, len(tagIds))
	for i := range tagIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	// SQL query to retrieve post IDs based on tag IDs
	query := fmt.Sprintf(`
		SELECT DISTINCT posttags_postid
		FROM cvwo_assignment.posttags
		WHERE posttags_tagid IN (%s)
	`, strings.Join(placeholders, ","))
	// Convert tag IDs to interface{} slice for variadic parameters
	tagIDArgs := make([]interface{}, len(tagIds))
	for i, id := range tagIds {
		tagIDArgs[i] = id
	}
	rows, err := models.DataBase.Query(query, tagIDArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var postIds []int
	for rows.Next() {
		var postId int
		err := rows.Scan(&postId)
		if err != nil {
			return nil, err
		}
		postIds = append(postIds, postId)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	fmt.Println("Post ids:", postIds)
	return postIds, nil

	//worked with psotgres
}

func getFilteredPosts(tagNames []string, sortBy string, pageStr string, perPageStr string) ([]models.Post, error) {
	tagIDs, err := getTagIdFromTagName(tagNames)
	if err != nil {
		fmt.Println("Get tagIds error:", err)
		return nil, err
	}
	postIDs, err := getPostIdsFromTagId(tagIDs)
	if err != nil {
		fmt.Println("Get postIDs error:", err)
		return nil, err
	}
	placeholders := make([]string, len(postIDs))
	if len(postIDs) > 0 {
		for i := range postIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
	} else {
		return make([]models.Post, 0), nil
	}

	// SQL query to retrieve posts based on post IDs
	query := fmt.Sprintf(`
		SELECT *
		FROM cvwo_assignment.posts
		WHERE postid IN (%s)
	`, strings.Join(placeholders, ","))
	//Sorting options
	switch sortBy {
	case "time":
		query += " ORDER BY datetime DESC"
	case "upvotes":
		query += " ORDER BY upvotes DESC"
	case "postid":
		query += " ORDER BY postid DESC"
	default:
		query += " ORDER BY postid DESC"
	}
	// Convert post IDs to interface{} slice for variadic parameters
	postIDArgs := make([]interface{}, len(postIDs))
	for i, id := range postIDs {
		postIDArgs[i] = id
	}

	if pageStr != "" && perPageStr != "" {
		pageNumber, _ := strconv.Atoi(pageStr)
		itemsPerPage, _ := strconv.Atoi(perPageStr)

		if pageNumber <= 0 {
			pageNumber = 1
		}

		if itemsPerPage <= 0 {
			itemsPerPage = 5
		}

		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(postIDArgs)+1, len(postIDArgs)+2)
		postIDArgs = append(postIDArgs, itemsPerPage, (pageNumber-1)*itemsPerPage)
	}

	rows, err := models.DataBase.Query(query, postIDArgs...)
	fmt.Println(query)
	fmt.Println(postIDArgs...)
	if err != nil {
		fmt.Println("Error querying:", err)
		return nil, err
	}
	defer rows.Close()
	posts, err := scanPosts(rows)
	if err != nil {
		fmt.Println("Error scanning:", err)
		return nil, err
	}
	return posts, nil
}
