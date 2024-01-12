package models

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	_ "github.com/lib/pq"
)

var DataBase *sql.DB

func InitDB() {
	//Opens database connection and creates the neccessary tables.
	// dbUser := os.Getenv("DB_USER")
	// dbPass := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")

	// Create a Postgres sql DSN (Data Source Name)
	// dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPass, dbName, dbHost, dbPort)

	dbURL := os.Getenv("DB_URL")
	var err error
	DataBase, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = DataBase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PostgresSql!")

	createUserTable()
	createPostsTable()
	createCommentsTable()
	createTagsTable()
	createUpvotesTable()
}

func createUserTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.users(
		username VARCHAR(50) NOT NULL,
		password VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (username)
	);`

	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created user table in my sql")
}

func createPostsTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.posts (
		postid SERIAL PRIMARY KEY,
		username VARCHAR(50) REFERENCES cvwo_assignment.users(username) ON DELETE CASCADE ON UPDATE CASCADE,
		tag JSON,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		upvotes INT CHECK (upvotes >= 0),
		datetime TIMESTAMP
	);`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created posts table in my sql")
}

func createCommentsTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.comments (
		commentid SERIAL PRIMARY KEY,
		comments_postid INT NOT NULL,
		comments_username VARCHAR(50) NOT NULL,
		description TEXT,
		upvotes INT CHECK (upvotes >= 0),
		datetime TIMESTAMP,
		CONSTRAINT comments_postid
			FOREIGN KEY (comments_postid)
			REFERENCES cvwo_assignment.posts (postid)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		CONSTRAINT comments_username
			FOREIGN KEY (comments_username)
			REFERENCES cvwo_assignment.users (username)
			ON DELETE NO ACTION
			ON UPDATE NO ACTION
	);`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created comments table in mysql")
}

func createTagsTable() {
	// clearTable := `	DROP TABLE IF EXISTS cvwo_assignment.tags;`
	// if _, err := DataBase.Exec(clearTable); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Deleted potential tags table")
	//Create table to store tags
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.tags (
		tagsid SERIAL PRIMARY KEY,
		name VARCHAR(80) UNIQUE
	);`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created tags table in mysql")

	// insertTags := `
	// INSERT INTO cvwo_assignment.tags (name)
	// VALUES
	// 	('recommendation'),
	// 	('compatibility'),
	// 	('troubleshooting'),
	// 	('deals');
	// `
	// if _, err := DataBase.Exec(insertTags); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Inserted tags into tags table")

	//Create many to many table to filter posts by tags
	createPostTagsTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.posttags (
		id SERIAL PRIMARY KEY,
		posttags_tagid INT NOT NULL,
		posttags_postid INT,
		CONSTRAINT postid
			FOREIGN KEY (posttags_postid)
			REFERENCES cvwo_assignment.posts (postid)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		CONSTRAINT tagid
			FOREIGN KEY (posttags_tagid)
			REFERENCES cvwo_assignment.tags (tagsid)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);`
	if _, err := DataBase.Exec(createPostTagsTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created posttags table in mysql")
}

func createUpvotesTable() {
	//Creates many to many upvotes table
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.upvotes (
		id SERIAL PRIMARY KEY,
		upvotes_postid INT,
		upvotes_commentid INT,
		upvotes_username VARCHAR(80),
		datetime TIMESTAMP,
		CONSTRAINT upvotes_postid
			FOREIGN KEY (upvotes_postid)
			REFERENCES cvwo_assignment.posts (postid)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		CONSTRAINT upvotes_commentid
			FOREIGN KEY (upvotes_commentid)
			REFERENCES cvwo_assignment.comments (commentid)
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		CONSTRAINT upvotes_username
			FOREIGN KEY (upvotes_username)
			REFERENCES cvwo_assignment.users (username)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
	`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created upvotes table in mysql")
}
