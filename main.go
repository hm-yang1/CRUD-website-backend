package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/models"
	"server/router"
	"server/sessions"
	// "github.com/joho/godotenv"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	sessions.InitSession()
	models.InitDB()
	r := router.Router()
	fmt.Println("Server running on port 8080")
	handler := router.SetupCORS(r)
	http.Handle("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
