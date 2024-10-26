package service

import (
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type CurrencyService struct {
	repo repository.Currency
	i18n config.I18nConfig
}

func NewCurrencyService(repo repository.Currency, i18n config.I18nConfig) *CurrencyService {
	return &CurrencyService{repo: repo, i18n: i18n}
}

func (s *CurrencyService) CreateCurrency(userID string, data *domain.CurrencyInput) (domain.Currency, error) {
	return s.repo.CreateCurrency(userID, data)
}

func (s *CurrencyService) GetCurrency(id string) (domain.Currency, error) {
	return s.repo.GetCurrency(id)
}

func (s *CurrencyService) FindCurrency(params domain.RequestParams) (domain.Response[domain.Currency], error) {
	return s.repo.FindCurrency(params)
}

func (s *CurrencyService) UpdateCurrency(id string, data interface{}) (domain.Currency, error) {
	return s.repo.UpdateCurrency(id, data)
}
func (s *CurrencyService) DeleteCurrency(id string) (domain.Currency, error) {
	return s.repo.DeleteCurrency(id)
}
