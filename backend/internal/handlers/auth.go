package handlers

import (
	"encoding/json"
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
