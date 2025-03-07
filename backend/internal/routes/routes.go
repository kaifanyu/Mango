package routes

import (
	"mango/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetUpRoutes() http.Handler {
	router := mux.NewRouter()

	// Authentication
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")
	router.HandleFunc("/signup", handlers.SignUpHandler).Methods("POST")
	router.HandleFunc("/auth/check", handlers.AuthCheck).Methods("GET", "POST")

	// Content Hosting
	router.HandleFunc("/api/content/audiobooks", handlers.AudioBookHandler).Methods("GET")
	router.HandleFunc("/api/content/anime", handlers.AnimeHandler).Methods("GET")
	router.HandleFunc("/api/content/movies", handlers.MovieHandler).Methods("GET")
	router.HandleFunc("/api/content/mangas", handlers.MangaHandler).Methods("GET")

	// User profile
	router.HandleFunc("/user/profile", handlers.GetUserProfileHandler).Methods("GET")

	staticDir := "../static" // Go up one level from cmd to backend, then into static
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	return c.Handler(router)
}
