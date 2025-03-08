package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mango/internal/database"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type User struct {
	Username string `json:"username"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	password := r.FormValue("password")

	// validate the user
	if !database.ValidateUser(username, password) {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// create a jwt claim that will expire in a day
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// create new token with HS256 algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// set auth_token to login
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false, //change to true in prod
		Path:     "/",
	})

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	// max 10 mb
	r.ParseMultipartForm(10 << 20)

	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	// gets image
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file"+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		http.Error(w, "Failed to Hash Image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	imageHash := hex.EncodeToString(hasher.Sum(nil))

	// verify registration image
	if !database.ValidateImageHash(imageHash) {
		http.Error(w, "Invalid registration image", http.StatusUnauthorized)
		return
	}

	// if it works
	database.RegisterUser(username, email, password)

	json.NewEncoder(w).Encode(map[string]string{"message": "Sign up successful"})
}

func AuthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("Verifying user")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println("cookie: ", cookie)

	tokenString := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user := User{Username: claims.Username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUserNameHandler(r *http.Request) (string, error) {

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", errors.New("Invalid Cookie")
	}

	// Extract the token string
	tokenString := cookie.Value

	// Initialize claims struct
	claims := &Claims{}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("Unauthorized token")
	}

	return claims.Username, nil
}



// return username and profile picture (TODO)
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	username, err := GetUserNameHandler(r)
	if err != nil {
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"username": username})

}
