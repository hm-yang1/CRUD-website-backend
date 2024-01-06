package sessions

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
)

// Pls store the keys in env variables thks
var Store *sessions.CookieStore

var jwtKey []byte

func InitSession() {
	sessionPassword := os.Getenv("SESSION_SECRET")
	if sessionPassword == "" {
		panic("SESSION_SECRET environment variable not set")
	}
	jwtPassword := os.Getenv("JWT_SECRET")
	if jwtPassword == "" {
		panic("JWT_SECRET environment variable not set")
	}
	jwtKey = []byte(jwtPassword)
	Store = sessions.NewCookieStore([]byte(sessionPassword))
	fmt.Println("Session store created")
}

func CreateJWT(username string) (string, error) {

	claims := jwt.MapClaims{
		"username":   username,
		"expiration": time.Now().Add(time.Hour * 24 * 7).Unix(), //token expiration time,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
