package blogs

import (
	"github.com/Sekarfo/P_blog/models"
	"gorm.io/gorm"
)

type BlogService interface {
	CreateBlog(blog *models.Blog) (*models.Blog, error)
}

type service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) BlogService {
	return &service{db: db}
}

func (s *service) CreateBlog(blog *models.Blog) (*models.Blog, error) {
	if err := s.db.Create(blog).Error; err != nil {
		return nil, err
	}
	return blog, nil
}
