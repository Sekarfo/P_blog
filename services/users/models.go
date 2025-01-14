package users

import "github.com/Sekarfo/P_blog/models"

type UsersService interface {
	CreateUser(user *models.User) (*models.User, error)
	LoginUser(email, password string) (*models.User, error)
	GetByParams(params *SearchParams) ([]models.User, error)
	GetByID(userID int) (*models.User, error)
	DeleteUser(userID int) error
	UpdateUser(user *models.User) (*models.User, error)
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
	SortByAsc  SortBy = "asc"
	SortByDesc SortBy = "desc"
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
