package blogs

import (
	"encoding/json"
	"net/http"

	"github.com/Sekarfo/P_blog/models"
	"github.com/Sekarfo/P_blog/services/blogs"
)

type Controller struct {
	blogService blogs.BlogService
}

func NewController(blogService blogs.BlogService) *Controller {
	return &Controller{blogService: blogService}
}

func (c *Controller) CreateBlog(w http.ResponseWriter, r *http.Request) {
	var blog models.Blog
	if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdBlog, err := c.blogService.CreateBlog(&blog)
	if err != nil {
		http.Error(w, "Error creating blog", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBlog)
}
