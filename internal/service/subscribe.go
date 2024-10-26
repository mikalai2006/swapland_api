package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type SubscribeService struct {
	repo repository.Subscribe
}

func NewSubscribeService(repo repository.Subscribe) *SubscribeService {
	return &SubscribeService{repo: repo}
}

func (s *SubscribeService) FindSubscribe(params domain.RequestParams) (domain.Response[model.Subscribe], error) {
	return s.repo.FindSubscribe(params)
}

func (s *SubscribeService) CreateSubscribe(userID string, data *model.SubscribeInput) (*model.Subscribe, error) {
	return s.repo.CreateSubscribe(userID, data)
}

func (s *SubscribeService) UpdateSubscribe(id string, userID string, Subscribe *model.Subscribe) (*model.Subscribe, error) {
	return s.repo.UpdateSubscribe(id, userID, Subscribe)
}

func (s *SubscribeService) DeleteSubscribe(id string) (model.Subscribe, error) {
	return s.repo.DeleteSubscribe(id)
}
