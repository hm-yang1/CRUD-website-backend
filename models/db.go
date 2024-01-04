package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DataBase *sql.DB

func InitDB() {
	//pls use env variables for the database
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Create a MySQL DSN (Data Source Name)
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Open a database connection
	var err error
	DataBase, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection
	err = DataBase.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL!")

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
		postid INT NOT NULL AUTO_INCREMENT,
		username VARCHAR(50) NULL,
		tag JSON NULL,
		title MEDIUMTEXT NOT NULL,
		description LONGTEXT NOT NULL,
		upvotes INT ZEROFILL NULL CHECK (upvotes >= 0),
		datetime TIMESTAMP NULL,
		PRIMARY KEY (postid),
		UNIQUE INDEX postid_UNIQUE (postid ASC) VISIBLE,
		INDEX username_idx (username ASC) VISIBLE,
		CONSTRAINT username
			FOREIGN KEY (username)
			REFERENCES cvwo_assignment.users (username)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created posts table in my sql")
}

func createCommentsTable() {
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.comments(
		commentid INT NOT NULL AUTO_INCREMENT,
		comments_postid INT NOT NULL,
		comments_username VARCHAR(50) NOT NULL,
		description LONGTEXT NULL,
		upvotes INT ZEROFILL NULL CHECK(upvotes >= 0),
		datetime TIMESTAMP NULL,
		PRIMARY KEY (commentid),
		UNIQUE INDEX commentid_UNIQUE (commentid ASC) VISIBLE,
		INDEX comments_username_idx (comments_username ASC) VISIBLE,
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
	fmt.Println("Deleted potential tags table")
	//Create table to store tags
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.tags (
		tagsid INT NOT NULL AUTO_INCREMENT,
		name VARCHAR(80) NULL,
		PRIMARY KEY (tagsid),
		UNIQUE INDEX tagsid_UNIQUE (tagsid ASC) VISIBLE
	);`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created tags table in mysql")

	// insertTags := `
	// INSERT INTO
	// 	cvwo_assignment.tags(name)
	// VALUES
	// 	('recommendation'),
	// 	('compatibility'),
	// 	('troubleshooting'),
	// 	('deals');
	// `
	// if _, err := DataBase.Exec(insertTags); err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println("Inserted tags into tags table")

	//Create many to many table to filter posts by tags
	createPostTagsTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.posttags (
		id INT NOT NULL AUTO_INCREMENT,
		posttags_tagid INT NOT NULL,
		posttags_postid INT NULL,
		PRIMARY KEY (id),
		UNIQUE INDEX id_UNIQUE (id ASC) VISIBLE,
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
	createTable := `
	CREATE TABLE IF NOT EXISTS cvwo_assignment.upvotes (
		id INT NOT NULL AUTO_INCREMENT,
		upvotes_postid INT NULL,
		upvotes_commentid INT NULL,
		upvotes_username VARCHAR(80) NULL,
		datetime TIMESTAMP NULL,
		PRIMARY KEY (id),
		UNIQUE INDEX id_UNIQUE (id ASC) VISIBLE,
		INDEX commentid_idx (upvotes_commentid ASC) VISIBLE,
		INDEX username_idx (upvotes_username ASC) VISIBLE,
		INDEX postid_idx (upvotes_postid ASC) VISIBLE,
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
	);`
	if _, err := DataBase.Exec(createTable); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created upvotes table in mysql")
}
