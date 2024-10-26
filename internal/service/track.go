package service

import (
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
)

type TrackService struct {
	repo repository.Track
}

func NewTrackService(repo repository.Track) *TrackService {
	return &TrackService{repo: repo}
}

func (s *TrackService) FindTrack(params domain.RequestParams) (domain.Response[domain.Track], error) {
	return s.repo.FindTrack(params)
}

func (s *TrackService) GetAllTrack(params domain.RequestParams) (domain.Response[domain.Track], error) {
	return s.repo.GetAllTrack(params)
}

func (s *TrackService) CreateTrack(userID string, track *domain.Track) (*domain.Track, error) {
	return s.repo.CreateTrack(userID, track)
}
