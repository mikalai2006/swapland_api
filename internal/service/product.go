package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type ProductService struct {
	repo        repository.Product
	userService *UserService
	Hub         *Hub
}

func NewProductService(repo repository.Product, userService *UserService, hub *Hub) *ProductService {
	return &ProductService{repo: repo, userService: userService, Hub: hub}
}

func (s *ProductService) FindProduct(params *model.ProductFilter) (domain.Response[model.Product], error) {
	return s.repo.FindProduct(params)
}

func (s *ProductService) CreateProduct(userID string, node *model.ProductInputData) (*model.Product, error) {
	result, err := s.repo.CreateProduct(userID, node)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(userID, model.UserStat{AddProduct: 1})
	}

	return result, err
}

func (s *ProductService) UpdateProduct(id string, userID string, data *model.Product) (*model.Product, error) {
	product, err := s.repo.UpdateProduct(id, userID, data)

	s.Hub.HandleMessage(domain.Message{Type: "message", Method: "update", Sender: userID, Recipient: "user2", Content: product, ID: "room1", Service: "product"})

	return product, err
}

func (s *ProductService) DeleteProduct(id string) (model.Product, error) {
	result, err := s.repo.DeleteProduct(id)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(result.UserID.Hex(), model.UserStat{AddProduct: -1})
	}

	return result, err
}
