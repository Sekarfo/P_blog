package users

import (
	"encoding/json"
	"net/http"
	"personal_blog/models"
	"personal_blog/services/users"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct {
	usersService users.UsersService
}

func NewController(usersS users.UsersService) *Controller {
	return &Controller{usersService: usersS}
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Send to service
	createdUser, err := c.usersService.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (c *Controller) GetByParams(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	params := users.SearchParams{}
	if r.URL.Query().Has("name") {
		name := r.URL.Query().Get("name")
		params.Name = &name
	}
	if r.URL.Query().Has("email") {
		email := r.URL.Query().Get("email")
		params.Email = &email
	}
	if r.URL.Query().Has("age") {
		age, err := strconv.Atoi(r.URL.Query().Get("age"))
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		params.Age = &age
	}
	if r.URL.Query().Has("sortBy") {
		sortBy := users.SortByFromString(
			r.URL.Query().Get("sortBy"),
		)
		params.SortBy = &sortBy
	}
	if r.URL.Query().Has("limit") {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		params.Limit = &limit
	}
	if r.URL.Query().Has("offset") {
		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		params.Offset = &offset
	}

	// Send to service
	users, err := c.usersService.GetByParams(&params)
	if err != nil {
		http.Error(w, "Error getting users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (c *Controller) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Send to service
	user, err := c.usersService.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user := models.User{
		Name:     r.FormValue("name"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Send to service
	updatedUser, err := c.usersService.UpdateUser(&user, id)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedUser)
}

func (c *Controller) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Send to service
	if err := c.usersService.DeleteUser(id); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
