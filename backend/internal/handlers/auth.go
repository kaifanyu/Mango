package handlers

import (
	"encoding/json"
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
	// create credintials format
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// decode the response
	json.NewDecoder(r.Body).Decode(&creds)

	// validate the user
	if !database.ValidateUser(creds.Username, creds.Password) {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	// create a jwt claim that will expire in a day
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := Claims{
		Username: creds.Username,
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

	tokenString := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user := User{Username: claims.Username}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
