package models

import "errors"

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
}

type ArticleSearch struct {
	Query    string             `json:"q"`
	SortBy   ArticleSeachSortBy `json:"sortBy"`
	Language string             `json:"language"`
	PageSize int                `json:"pageSize"`
	Page     int                `json:"page"`
}

type ArticleSeachSortBy string

const (
	ArticleSeachSortByRelevancy   = "relevancy"
	ArticleSeachSortByPopularity  = "popularity"
	ArticleSeachSortByPublishedAt = "publishedAt"
)

func (aBy ArticleSeachSortBy) String() string {
	return string(aBy)
}

func ArticleSeachSortByFromString(s string) (ArticleSeachSortBy, error) {
	switch s {
	case "relevancy":
		return ArticleSeachSortByRelevancy, nil
	case "popularity":
		return ArticleSeachSortByPopularity, nil
	case "publishedAt":
		return ArticleSeachSortByPublishedAt, nil
	default:
		return ArticleSeachSortByRelevancy, errors.New("Invalid sort by")
	}
}
