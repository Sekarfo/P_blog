package articles

import (
	"net/url"
	"personal_blog/models"
	"strconv"
)

type toNewsAPIReq struct {
	models.ArticleSearch
	ApiKey string `json:"apiKey"`
}

func (r *toNewsAPIReq) toQueryParams() string {
	values := url.Values{}
	values.Set("q", r.Query)
	values.Set("sortBy", r.SortBy.String())
	values.Set("pageSize", strconv.Itoa(r.PageSize))
	values.Set("page", strconv.Itoa(r.Page))
	values.Set("apiKey", r.ApiKey)

	return "?" + values.Encode()
}

type newsResponse struct {
	Status       string           `json:"status"`
	TotalResults int              `json:"totalResults"`
	Articles     []models.Article `json:"articles"`
}
