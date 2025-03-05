package routes

import (
	"mango/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetUpRoutes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	router.HandleFunc("/auth/check", handlers.AuthCheck).Methods("GET", "POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(router)
}
