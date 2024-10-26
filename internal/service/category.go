package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type CategoryService struct {
	repo repository.Category
	i18n config.I18nConfig
}

func NewCategoryService(repo repository.Category, i18n config.I18nConfig) *CategoryService {
	return &CategoryService{repo: repo, i18n: i18n}
}

func (s *CategoryService) FindCategory(params domain.RequestParams) (domain.Response[model.Category], error) {
	return s.repo.FindCategory(params)
}

func (s *CategoryService) GetAllCategory(params domain.RequestParams) (domain.Response[model.Category], error) {
	return s.repo.GetAllCategory(params)
}

func (s *CategoryService) CreateCategory(userID string, Category *model.Category) (*model.Category, error) {
	return s.repo.CreateCategory(userID, Category)
}

func (s *CategoryService) UpdateCategory(id string, userID string, Category *model.CategoryInput) (*model.Category, error) {
	return s.repo.UpdateCategory(id, userID, Category)
}

func (s *CategoryService) DeleteCategory(id string) (model.Category, error) {
	return s.repo.DeleteCategory(id)
}
