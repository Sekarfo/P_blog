package users

import (
	"errors"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/Sekarfo/P_blog/utils"

	"github.com/Sekarfo/P_blog/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) UsersService {
	return &service{db: db}
}

func (s *service) CreateUser(user *models.User) (*models.User, error) {
	if err := utils.ValidateUserInput(user); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error hashing password")
	}
	user.Password = string(hashedPassword)

	createDBError := s.db.Create(user).Error
	if createDBError != nil {
		return nil, createDBError
	}
	return user, nil
}

func (s *service) GetByParams(params *SearchParams) ([]models.User, error) {
	var users []models.User
	query := s.db

	if params.Name != nil {
		query = query.Where("name = ?", *params.Name)
	}
	if params.Email != nil {
		query = query.Where("email = ?", *params.Email)
	}
	if params.Age != nil {
		query = query.Where("age = ?", *params.Age)
	}
	if params.SortBy != nil {
		query = query.Order("name " + string(*params.SortBy))
	}
	if params.Limit != nil {
		query = query.Limit(*params.Limit)
	}
	if params.Offset != nil {
		query = query.Offset(*params.Offset)
	}

	query.Find(&users)
	if query.Error != nil {
		return nil, query.Error
	}
	return users, nil
}

func (s *service) GetByID(userID int) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *service) UpdateUser(user *models.User) (*models.User, error) {
	existingUser, err := s.GetByID(int(user.ID))
	if err != nil {
		return nil, err
	}

	if user.Name != "" {
		existingUser.Name = user.Name
	}
	if user.Email != "" {
		existingUser.Email = user.Email
	}

	if err := s.db.Save(existingUser).Error; err != nil {
		return nil, err
	}

	return existingUser, nil
}

func (s *service) DeleteUser(userID int) error {
	if err := s.db.Delete(&models.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (s *service) LoginUser(email, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}

func SendSupportEmail(userID int, message string, handler *multipart.FileHeader, file multipart.File) error {
	// Load email configuration from environment variables
	fromEmail := os.Getenv("EMAIL_FROM")
	toEmail := os.Getenv("EMAIL_TO")
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Create email message
	email := gomail.NewMessage()
	email.SetHeader("From", fromEmail)
	email.SetHeader("To", toEmail)
	email.SetHeader("Subject", fmt.Sprintf("Support Request from User ID: %d", userID))
	email.SetBody("text/plain", message)

	// Attach file if provided
	if handler != nil && file != nil {
		email.Attach(handler.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, file)
			return err
		}))
	}

	// Send email
	dialer := gomail.NewDialer(smtpServer, smtpPort, smtpUser, smtpPassword)
	if err := dialer.DialAndSend(email); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
