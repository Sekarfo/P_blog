package users

import (
	"encoding/json"
	"github.com/Sekarfo/P_blog/models"
	"github.com/Sekarfo/P_blog/services/users"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type Controller struct {
	sessionsStore *sessions.CookieStore
	usersService  users.UsersService
}

func NewController(
	usersS users.UsersService,
) *Controller {
	storage := sessions.NewCookieStore([]byte(os.Getenv("SESSIONS_KEY")))
	return &Controller{
		sessionsStore: storage,
		usersService:  usersS,
	}
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

func (c *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	log.Println("Attempting login with email:", email)

	// Send to service
	user, err := c.usersService.LoginUser(email, password)
	if err != nil {
		log.Println("Invalid login attempt for email:", email)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	log.Println("Login successful for email:", email)

	// Create session
	session := sessions.NewSession(c.sessionsStore, "session-name")
	session.Values["user_id"] = user.ID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email
	err = session.Save(r, w)

	if err != nil {
		log.Println("Error saving session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Login successful",
	})
}

func (c *Controller) GetProfile(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionsStore.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := c.usersService.GetByID(int(userID))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Exclude sensitive information
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (c *Controller) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	session, err := c.sessionsStore.Get(r, "session-name")
	if err != nil {
		log.Println("Error getting session:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, ok := session.Values["user_id"].(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user.ID = userID

	updatedUser, err := c.usersService.UpdateUser(&user)
	if err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (c *Controller) SendSupportRequest(w http.ResponseWriter, r *http.Request) {
	// Get user session
	session, err := c.sessionsStore.Get(r, "session-name")
	if err != nil {
		log.Printf("Error getting session: %v\n", err)
		http.Error(w, "Failed to retrieve session. Please try again.", http.StatusInternalServerError)
		return
	}

	// Validate user ID from session
	userID, ok := session.Values["user_id"].(uint)
	if !ok {
		http.Error(w, "Unauthorized. Please log in.", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Invalid form submission. Ensure fields are correct.", http.StatusBadRequest)
		return
	}

	// Get form values
	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Message field cannot be empty.", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("attachment")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("Error reading file: %v\n", err)
		http.Error(w, "Failed to read attachment.", http.StatusInternalServerError)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// Call email service to send the support request
	err = users.SendSupportEmail(int(userID), message, handler, file)
	if err != nil {
		log.Printf("Error sending email: %v\n", err)
		http.Error(w, "Failed to send support request. Please try again.", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Support request sent successfully."))
}
