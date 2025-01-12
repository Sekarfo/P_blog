package users

import "github.com/Sekarfo/P_blog/models"

type UsersService interface {
	CreateUser(user *models.User) (*models.User, error)
	GetByParams(params *SearchParams) ([]models.User, error)
	GetByID(userID int) (*models.User, error)
	UpdateUser(user *models.User, updatedUserID int) (*models.User, error)
	DeleteUser(userID int) error
}

type SearchParams struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Age    *int    `json:"age"`
	SortBy *SortBy `json:"sortBy"`
	Limit  *int    `json:"limit"`
	Offset *int    `json:"offset"`
}

type SortBy string

const (
	SortByAsc  = "asc"
	SortByDesc = "desc"
)

func SortByFromString(s string) SortBy {
	switch s {
	case "asc":
		return SortByAsc
	case "desc":
		return SortByDesc
	default:
		return SortByAsc
	}
}
