package comments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hitpads/reado_ap/internal/auth"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
)

func AddComment(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Content string `json:"content"`
		PostID  uint   `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(auth.UserIDKey).(uint)

	comment := models.Comment{
		Content:  payload.Content,
		PostID:   payload.PostID,
		AuthorID: userID,
	}

	if err := db.DB.Create(&comment).Error; err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	var comments []models.Comment
	if err := db.DB.Preload("Author").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	var response []map[string]interface{}
	for _, comment := range comments {
		response = append(response, map[string]interface{}{
			"id":      comment.ID,
			"content": comment.Content,
			"author":  comment.Author.Name,
			"created": comment.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
