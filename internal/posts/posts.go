package posts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hitpads/reado_ap/internal/auth"
	"github.com/hitpads/reado_ap/internal/db"
	"github.com/hitpads/reado_ap/internal/models"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	authorID := r.Context().Value(auth.UserIDKey).(uint)
	post.AuthorID = authorID

	if err := db.DB.Create(&post).Error; err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(post)
	if err != nil {
		return
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	if err := db.DB.Preload("Author").Find(&posts).Error; err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	var postResponses []map[string]interface{}
	for _, post := range posts {
		postResponses = append(postResponses, map[string]interface{}{
			"id":      post.ID,
			"title":   post.Title,
			"content": post.Content,
			"author":  post.Author.Name,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postResponses)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	if err := db.DB.Delete(&models.Post{}, postID).Error; err != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetPostDetails(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	var post models.Post

	if err := db.DB.Preload("Author").First(&post, postID).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"id":      post.ID,
		"title":   post.Title,
		"content": post.Content,
		"author":  post.Author.Name,
		"likes":   post.Likes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := r.Context().Value(auth.UserIDKey).(uint)

	// Check if the post exists
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check if the user has already liked the post
	var like models.Like
	if err := db.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err == nil {
		http.Error(w, "You have already liked this post", http.StatusBadRequest)
		return
	}

	// Add the like
	like = models.Like{
		PostID: post.ID,
		UserID: userID,
	}
	if err := db.DB.Create(&like).Error; err != nil {
		http.Error(w, "Failed to like post", http.StatusInternalServerError)
		return
	}

	// Increment the post's like count
	post.Likes++
	if err := db.DB.Save(&post).Error; err != nil {
		http.Error(w, "Failed to update like count", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likes": post.Likes,
	})
}

func UnlikePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := r.Context().Value(auth.UserIDKey).(uint)

	// Check if the post exists
	var post models.Post
	if err := db.DB.First(&post, postID).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Check if the user has liked the post
	var like models.Like
	if err := db.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		http.Error(w, "You have not liked this post", http.StatusBadRequest)
		return
	}

	// Remove the like
	if err := db.DB.Delete(&like).Error; err != nil {
		http.Error(w, "Failed to remove like", http.StatusInternalServerError)
		return
	}

	// Decrement the post's like count
	if post.Likes > 0 {
		post.Likes--
	}
	if err := db.DB.Save(&post).Error; err != nil {
		http.Error(w, "Failed to update like count", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"likes": post.Likes,
	})
}

func GetLikeStatus(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")
	userID := r.Context().Value(auth.UserIDKey).(uint)

	var like models.Like
	if err := db.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err != nil {
		json.NewEncoder(w).Encode(map[string]bool{"liked": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"liked": true})
}
