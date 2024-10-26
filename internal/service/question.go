package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type QuestionService struct {
	repo repository.Question
	Hub  *Hub
}

func NewQuestionService(repo repository.Question, hub *Hub) *QuestionService {
	return &QuestionService{repo: repo, Hub: hub}
}

func (s *QuestionService) FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error) {
	return s.repo.FindQuestion(params)
}

func (s *QuestionService) CreateQuestion(userID string, tag *model.QuestionInput) (*model.Question, error) {
	result, err := s.repo.CreateQuestion(userID, tag)

	s.Hub.HandleMessage(domain.Message{Type: "message", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "question"})

	return result, err
}

func (s *QuestionService) UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error) {
	result, err := s.repo.UpdateQuestion(id, userID, data)

	s.Hub.HandleMessage(domain.Message{Type: "message", Sender: userID, Recipient: "user2", Content: result, ID: "room1", Service: "question"})

	return result, err
}

func (s *QuestionService) DeleteQuestion(id string) (model.Question, error) {
	return s.repo.DeleteQuestion(id)
}
