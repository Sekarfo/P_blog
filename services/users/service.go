package users

import (
	"errors"

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

func (s *service) UpdateUser(user *models.User, _ int) (*models.User, error) {
	if err := utils.ValidateUserInput(user); err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Error hashing password")
	}
	user.Password = string(hashedPassword)

	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) DeleteUser(userID int) error {
	if err := s.db.Delete(&models.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}
