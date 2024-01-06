package models

import (
	"time"
)

type User struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Created_at time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Post struct {
	Postid      int64     `json:"postid"`
	Username    string    `json:"username"`
	Tags        []string  `json:"tags"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Upvote      int32     `json:"upvote"`
	Datetime    time.Time `json:"datetime"`
}

type PostRequest struct {
	Username    string   `json:"username"`
	Tags        []string `json:"tags"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
}

type Comment struct {
	Commentid   int64     `json:"commentid"`
	Postid      int64     `json:"postid"`
	Username    string    `json:"username"`
	Description string    `json:"description"`
	Upvote      int32     `json:"upvote"`
	Datetime    time.Time `json:"datetime"`
}

type CommentRequest struct {
	Postid      int64  `json:"postid"`
	Username    string `json:"username"`
	Description string `json:"description"`
}

type Tag struct {
	Tagid int8   `json:"tagid"`
	Name  string `json:"name"`
}

type Upvote struct {
	Postid    int64     `json:"postid"`
	Commentid int64     `json:"commentid"`
	Username  string    `json:"username"`
	Datetime  time.Time `json:"datetime"`
}

type UpvoteRequest struct {
	Postid    int64  `json:"postid"`
	Commentid int64  `json:"commentid"`
	Username  string `json:"username"`
}

type Sorting struct {
	Sortingid int8   `json:"sortingid"`
	Name      string `json:"name"`
}

func InitUser(username string, password string) User {
	u := User{
		Username:   username,
		Password:   password,
		Created_at: time.Now(),
	}
	return u
}

func InitPost(username string, tags []string, title string, descrption string) Post {
	p := Post{
		Username:    username,
		Tags:        tags,
		Title:       title,
		Description: descrption,
		Upvote:      0,
		Datetime:    time.Now(),
	}
	return p
}

func EditPost(postid int64, username string, tags []string, title string, descrption string) Post {
	p := Post{
		Postid:      postid,
		Username:    username,
		Tags:        tags,
		Title:       title,
		Description: descrption,
	}
	return p
}

func InitComment(postid int64, username string, descrption string) Comment {
	p := Comment{
		Postid:      postid,
		Username:    username,
		Description: descrption,
		Upvote:      0,
		Datetime:    time.Now(),
	}
	return p
}
