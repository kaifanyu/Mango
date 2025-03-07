package handlers

import (
	"encoding/json"
	"fmt"
	"mango/internal/database"
	"net/http"
)

var limit int = 20


type Anime struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	Author     string `json:"author"`
	Duration   int    `json:"duration"`
}

type Audiobook struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CoverImage string `json:"cover_image"`
	Author     string `json:"author"`
	Duration   int    `json:"duration"`
}



func AnimeHandler(w http.ResponseWriter, r*http.Request) {
	query := `
		SELECT c.id, c.title, a.episodes
		FROM content c
		JOIN anime a ON c.id = a.content_id
		WHERE c.content_type = 'anime'
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
