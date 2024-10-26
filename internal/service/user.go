package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type UserService struct {
	repo repository.User
	Hub  *Hub
}

func NewUserService(repo repository.User, hub *Hub) *UserService {
	return &UserService{repo: repo, Hub: hub}
}

func (s *UserService) GetUser(id string) (model.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) FindUser(params domain.RequestParams) (domain.Response[model.User], error) {
	return s.repo.FindUser(params)
}

func (s *UserService) CreateUser(userID string, user *model.User) (*model.User, error) {
	return s.repo.CreateUser(userID, user)
}

func (s *UserService) DeleteUser(id string) (model.User, error) {
	return s.repo.DeleteUser(id)
}

func (s *UserService) UpdateUser(id string, user *model.User) (model.User, error) {
	result, err := s.repo.UpdateUser(id, user)
	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "PATCH", Sender: id, Recipient: "user2", Content: result, ID: "room1", Service: "user"})

	return result, err
}

func (s *UserService) Iam(userID string) (model.User, error) {
	user, err := s.repo.Iam(userID)
	if err != nil {
		return user, err
	}

	// user, err = s.UpdateUser(userID, &model.User{Online: true})
	// s.Hub.HandleMessage(domain.Message{Type: "message", Sender: "user1", Recipient: "user2", Content: user, ID: "room1", Service: "user"})

	return user, err
}

func (s *UserService) SetStat(userID string, statData model.UserStat) (model.User, error) {
	return s.repo.SetStat(userID, statData)
}

func (s *UserService) SetBal(userID string, value int) (model.User, error) {
	return s.repo.SetBal(userID, value)
}
