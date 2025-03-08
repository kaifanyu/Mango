package handlers

import (
	"encoding/json"
	"mango/internal/database"
	"net/http"

	"github.com/gorilla/mux"
)

type ProgressRequest struct {
	VideoID  string `json:"video_id"`
	Progress int    `json:"progress"`
}

func SaveProgressHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username, err := GetUserNameHandler(r)
	if err != nil {
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}

	// get the userid
	userID, err := database.GetUserID(username)

	// save progress to database with video_id and progress
	var req ProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.VideoID == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	if req.Progress < 0 {
		http.Error(w, "Progress cannot be negative", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO user_progress (user_id, video_id, progress)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			progress = VALUES(progress),
			updated_at = CURRENT_TIMESTAMP
	`

	_, err = database.DB.Exec(query, userID, req.VideoID, req.Progress)
	if err != nil {
		http.Error(w, "Failed to save to database", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"progress": req.Progress,
		"video_id": req.VideoID,
	})

}

// Gets current user's video progress from database
func GetProgressHandler(w http.ResponseWriter, r *http.Request) {
	// returning json file
	w.Header().Set("Content-Type", "application/json")

	// first get the username
	username, err := GetUserNameHandler(r)
	if err != nil {
		http.Error(w, "Failed to retrieve username", http.StatusInternalServerError)
		return
	}

	// get the userid
	userID, err := database.GetUserID(username)

	// get the videoid
	vars := mux.Vars(r)
	videoID := vars["id"]

	// Find the userprogress
	query := `
		SELECT progress 
		FROM user_progress
		WHERE user_id = ? AND video_id = ?
	`

	results := database.ExecDisplayQuery(query, userID, videoID)

	if len(results) == 0 {
		json.NewEncoder(w).Encode(map[string]int{"progress": 0})
		return
	}

	progressData := results[0]
	progress := getIntValue(progressData["progress"])

	json.NewEncoder(w).Encode(map[string]int{"progress": progress})
}
