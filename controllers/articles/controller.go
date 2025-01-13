package articles

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Sekarfo/P_blog/models"

	articlesS "github.com/Sekarfo/P_blog/services/articles"
)

type Controller struct {
	articleService articlesS.ArticleService
}

func NewController(articleService articlesS.ArticleService) *Controller {
	return &Controller{articleService: articleService}
}

func (c *Controller) FetchArticles(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	params := models.ArticleSearch{}
	if r.URL.Query().Has("q") {
		query := r.URL.Query().Get("q")
		params.Query = query
	}
	if r.URL.Query().Has("sortBy") {
		sortBy, err := models.ArticleSeachSortByFromString(r.URL.Query().Get("sortBy"))
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		params.SortBy = sortBy
	}
	if r.URL.Query().Has("pageSize") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		params.PageSize = pageSize
	}
	if r.URL.Query().Has("page") {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		params.Page = page
	}
	if r.URL.Query().Has("language") {
		language := r.URL.Query().Get("language")
		params.Language = language
	}

	articles, total, err := c.articleService.GetArticles(params)
	if err != nil {
		http.Error(w, "Error fetching articles", http.StatusInternalServerError)
		return
	}
	out := outResponse{
		Status:   "success",
		Total:    total,
		Articles: articles,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		log.Print(fmt.Sprintf("Error encoding articles", err))
	}
}

type outResponse struct {
	Status   string           `json:"status"`
	Total    int              `json:"totalResults"`
	Articles []models.Article `json:"articles"`
}
