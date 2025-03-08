package handlers

import (
	"encoding/json"
	"fmt"
	"mango/internal/database"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var limit int = 20

type Anime struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	Episodes   int    `json:"episodes"`
}

type Audiobook struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	Author     string `json:"author"`
	Duration   int    `json:"duration"`
}

type Manga struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	Author     string `json:"author"`
	Chapters   int    `json:"Chapters"`
}

type Movie struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	Duration   int    `json:"duration"`
}

func MangaHandler(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT c.id, c.title, a.author, a.chapters
        FROM content c
        JOIN mangas a ON c.id = a.content_id
        WHERE c.content_type = 'manga'
        ORDER BY c.id DESC
        LIMIT ?
    `
	mangas := database.ExecDisplayQuery(query, limit)

	results := make([]Manga, 0)

	for _, manga := range mangas {
		id := getIntValue(manga["id"])
		coverImage := fmt.Sprintf("http://localhost:8080/static/content/covers/mangas/%d.jpg", id)

		results = append(results, Manga{
			ID:         id,
			Title:      getStringValue(manga["title"]),
			CoverImage: coverImage,
			Author:     getStringValue(manga["author"]),
			Chapters:   getIntValue(manga["duration"]),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func MovieHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT c.id, c.title, a.duration
		FROM content c
		JOIN movies a ON c.id = a.content_id
		WHERE c.content_type = 'movies'
		ORDER BY c.id DESC
		LIMIT ?
	`

	movies := database.ExecDisplayQuery(query, limit)

	results := make([]Movie, 0)

	for _, movie := range movies {
		id := getIntValue(movie["id"])
		coverImage := fmt.Sprintf("http://localhost:8080/static/content/covers/movies/%d.jpg", id)

		results = append(results, Movie{
			ID:         id,
			Title:      getStringValue(movie["title"]),
			CoverImage: coverImage,
			Duration:   getIntValue(movie["duration"]),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func AnimeHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT c.id, c.title, a.episodes
		FROM content c
		JOIN anime a ON c.id = a.content_id
		WHERE c.content_type = 'anime'
		ORDER BY c.id DESC
		LIMIT ?
	`

	animes := database.ExecDisplayQuery(query, limit)

	results := make([]Anime, 0)

	for _, anime := range animes {
		id := getIntValue(anime["id"])
		coverImage := fmt.Sprintf("http://localhost:8080/static/content/covers/anime/%d.jpg", id)

		results = append(results, Anime{
			ID:         id,
			Title:      getStringValue(anime["title"]),
			CoverImage: coverImage,
			Episodes:   getIntValue(anime["episodes"]),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)

}

func AudioBookHandler(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT c.id, c.title, a.author, a.duration
        FROM content c
        JOIN audiobooks a ON c.id = a.content_id
        WHERE c.content_type = 'audiobooks'
        ORDER BY c.id DESC
        LIMIT ?
    `
	audiobooks := database.ExecDisplayQuery(query, limit)

	results := make([]Audiobook, 0)

	for _, book := range audiobooks {
		id := getIntValue(book["id"])
		coverImage := fmt.Sprintf("http://localhost:8080/static/content/covers/audiobooks/%d.jpg", id)

		results = append(results, Audiobook{
			ID:         id,
			Title:      getStringValue(book["title"]),
			CoverImage: coverImage,
			Author:     getStringValue(book["author"]),
			Duration:   getIntValue(book["duration"]),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func AudiobookDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) //gets audibook id
	audiobookID := vars["id"]

	id, err := strconv.Atoi(audiobookID)
	if err != nil {
		http.Error(w, "Invalid audibook ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT c.id, c.title, a.author, a.duration
		FROM content c
		JOIN audiobooks a ON c.id = a.content_id
		WHERE c.id = ? and c.content_type = 'audiobooks'
	`

	results := database.ExecDisplayQuery(query, id)

	if len(results) == 0 {
		http.Error(w, "Audiobook not found", http.StatusNotFound)
		return
	}

	book := results[0]
	coverImage := fmt.Sprintf("http://localhost:8080/static/content/covers/audiobooks/%s.jpg", audiobookID)

	detail := Audiobook{
		ID:         getIntValue(audiobookID),
		Title:      getStringValue(book["title"]),
		CoverImage: coverImage,
		Author:     getStringValue(book["author"]),
		Duration:   getIntValue(book["duration"]),

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

// Helper functions for type assertions
func getStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

func getIntValue(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}
