package routes

import (
	"net/http"

	"mango/internal/handlers"

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

	// NEW: Content Detail APIs
	router.HandleFunc("/api/content/audiobooks/{id}", handlers.AudiobookDetailHandler).Methods("GET")

	// Video Progress
	router.HandleFunc("/api/progress/{id}", handlers.GetProgressHandler).Methods("GET")
	router.HandleFunc("/api/progress", handlers.SaveProgressHandler).Methods("POST")
	// User profile
	router.HandleFunc("/user/profile", handlers.GetUserProfileHandler).Methods("GET")

	// Create a CORS handler with your settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://www.thejoyestboy.com"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "Range"},
		ExposedHeaders:   []string{"Content-Length", "Content-Range", "Accept-Ranges"},
		AllowCredentials: true,
	})

	// Static files with CORS handling
	staticDir := "/media/para/APTX/static/"
	fs := http.FileServer(http.Dir(staticDir))

	// Create a CORS-enabled file server handler
	corsStaticHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers for static files
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // Dynamic origin handling
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Range, Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")

		// If it's an OPTIONS preflight request, just return OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, serve the file
		http.StripPrefix("/static/", fs).ServeHTTP(w, r)
	})

	// Use our custom handler for static files
	router.PathPrefix("/static/").Handler(corsStaticHandler)

	// Apply CORS middleware to the router for API routes
	return c.Handler(router)
}
