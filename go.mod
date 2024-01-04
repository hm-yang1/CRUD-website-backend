module server

go 1.21.5

require (
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/sessions v1.2.2
	golang.org/x/crypto v0.17.0
)

require github.com/joho/godotenv v1.5.1

require (
	// github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/rs/cors v1.10.1
)