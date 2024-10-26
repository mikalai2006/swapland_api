package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type NodeVoteService struct {
	repo repository.NodeVote
}

func NewNodeVoteService(repo repository.NodeVote) *NodeVoteService {
	return &NodeVoteService{repo: repo}
}

func (s *NodeVoteService) CreateNodeVote(userID string, data *model.NodeVote) (*model.NodeVote, error) {
	return s.repo.CreateNodeVote(userID, data)
}

func (s *NodeVoteService) FindNodeVote(params domain.RequestParams) (domain.Response[model.NodeVote], error) {
	return s.repo.FindNodeVote(params)
}

// func (s *NodeVoteService) GetAllNodeVote(params domain.RequestParams) (domain.Response[model.NodeVote], error) {
// 	return s.repo.GetAllNodeVote(params)
// }

func (s *NodeVoteService) UpdateNodeVote(id string, userID string, data *model.NodeVoteInput) (*model.NodeVote, error) {
	return s.repo.UpdateNodeVote(id, userID, data)
}

func (s *NodeVoteService) DeleteNodeVote(id string) (model.NodeVote, error) {
	return s.repo.DeleteNodeVote(id)
}
