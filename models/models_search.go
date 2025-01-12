package models

import "errors"

type ArticleSearch struct {
	Query    string             `json:"q"`
	SortBy   ArticleSeachSortBy `json:"sortBy"`
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
