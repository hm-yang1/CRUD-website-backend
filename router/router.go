package router

import (
	"fmt"
	"net/http"
	"server/handlers"
	"server/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	fmt.Println("Created mux router")

	//User authentication routes
	router.HandleFunc("/api/login", handlers.LoginHandler).Methods("Post")
	router.HandleFunc("/api/logout", middleware.AuthRequired(handlers.LogoutHandler)).Methods("Post")
	router.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/session", middleware.SessionHandler).Methods("GET")

	//Post routes
	router.HandleFunc("/api/posts", handlers.GetPostsHandler).Methods("GET")
	router.HandleFunc("/api/posts/{postid}", handlers.GetPostByIdHandler).Methods("GET")
	router.HandleFunc("/api/filtered/posts", handlers.GetFilteredPostsHandler).Methods("GET")
	router.HandleFunc("/api/search/posts", handlers.GetSearchedPostsHandler).Methods("GET")
	router.HandleFunc("/api/posts", middleware.AuthRequired(handlers.CreatePostHandler)).Methods("POST")
	router.HandleFunc("/api/posts/{postid}", middleware.IsPostEditAuthorized(handlers.EditPostHandler)).Methods("PUT")
	router.HandleFunc("/api/posts/{postid}", middleware.IsPostEditAuthorized(handlers.DeletePostHandler)).Methods("DELETE")

	//Comments routes
	router.HandleFunc("/api/posts/{postid}/comments", handlers.GetCommentsHandler).Methods("GET")
	router.HandleFunc("/api/comments/{commentid}", handlers.GetCommentHandler).Methods("GET")
	router.HandleFunc("/api/posts/{postid}/comments", middleware.AuthRequired(handlers.CreateCommentHandler)).Methods("POST")
	router.HandleFunc("/api/comments/{commentid}", middleware.IsCommentEditAuthorized(handlers.EditCommentHandler)).Methods("PUT")
	router.HandleFunc("/api/comments/{commentid}", middleware.IsCommentEditAuthorized(handlers.DeleteCommentHandler)).Methods("DELETE")

	//Upvotes routes
	router.HandleFunc("/api/upvotes/posts/{postid}", middleware.AuthRequired(handlers.GetPostUpvoteHandler)).Methods("GET")
	router.HandleFunc("/api/upvotes/comments/{commentid}", middleware.AuthRequired(handlers.GetCommentUpvoteHandler)).Methods("GET")
	router.HandleFunc("/api/upvotes/add", middleware.AuthRequired(handlers.AddUpvoteHandler)).Methods("POST")
	router.HandleFunc("/api/upvotes/remove", middleware.AuthRequired(handlers.DeleteUpvoteHandler)).Methods("POST")

	//Tags route
	router.HandleFunc("/api/tags", handlers.GetTagsHandler).Methods("GET")
	//Profile routes
	// router.HandleFunc("/api/profile/{username}", middleware.GetProfileHandler).Methods("GET")

	return router
}

func SetupCORS(handler http.Handler) http.Handler {
	// Create a CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://web-forum-yw1j.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Use the CORS handler
	return c.Handler(handler)
}
