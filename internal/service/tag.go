package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type TagService struct {
	repo repository.Tag
}

func NewTagService(repo repository.Tag) *TagService {
	return &TagService{repo: repo}
}

func (s *TagService) FindTag(params domain.RequestParams) (domain.Response[model.Tag], error) {
	return s.repo.FindTag(params)
}

func (s *TagService) GetAllTag(params domain.RequestParams) (domain.Response[model.Tag], error) {
	return s.repo.GetAllTag(params)
}

func (s *TagService) CreateTag(userID string, tag *model.Tag) (*model.Tag, error) {
	return s.repo.CreateTag(userID, tag)
}

func (s *TagService) UpdateTag(id string, userID string, tag *model.Tag) (*model.Tag, error) {
	return s.repo.UpdateTag(id, userID, tag)
}

func (s *TagService) DeleteTag(id string) (model.Tag, error) {
	return s.repo.DeleteTag(id)
}
