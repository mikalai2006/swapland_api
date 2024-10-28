package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type AddressService struct {
	repo repository.Address
	i18n config.I18nConfig
}

func NewAddressService(repo repository.Address, i18n config.I18nConfig) *AddressService {
	return &AddressService{repo: repo, i18n: i18n}
}

func (s *AddressService) CreateAddress(userID string, address domain.ResponseNominatim) (*model.Address, error) {
	return s.repo.CreateAddress(userID, address)
}

func (s *AddressService) FindAddress(input *model.AddressFilter) (domain.Response[model.Address], error) {
	return s.repo.FindAddress(input)
}

// func (s *AddressService) GetAllAddress(params domain.RequestParams) (domain.Response[model.Address], error) {
// 	return s.repo.GetAllAddress(params)
// }

func (s *AddressService) UpdateAddress(id string, userID string, data domain.ResponseNominatim) (*model.Address, error) {
	return s.repo.UpdateAddress(id, userID, data)
}

func (s *AddressService) DeleteAddress(id string) (model.Address, error) {
	return s.repo.DeleteAddress(id)
}

// func (s *PageService) GetPageForRouters() (domain.Response[domain.PageRoutes], error) {
// 	return s.repo.GetPageForRouters()
// }

// func (s *PageService) GetFullPage(params domain.RequestParams) (domain.Response[domain.Page], error) {
// 	return s.repo.GetFullPage(params)
// }

// func (s *PageService) UpdatePageWithContent(id string, data map[string]interface{}) (domain.Page, error) {
// 	return s.repo.UpdatePageWithContent(id, data)
// }
