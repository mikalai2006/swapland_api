package service

import (
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewService struct {
	repo        repository.Review
	userService *UserService
}

func NewReviewService(repo repository.Review, userService *UserService) *ReviewService {
	return &ReviewService{repo: repo, userService: userService}
}

func (s *ReviewService) FindReview(params domain.RequestParams) (domain.Response[model.Review], error) {
	return s.repo.FindReview(params)
}

func (s *ReviewService) GetAllReview(params domain.RequestParams) (domain.Response[model.Review], error) {
	return s.repo.GetAllReview(params)
}

func (s *ReviewService) CreateReview(userID string, review *model.Review) (*model.Review, error) {
	var result *model.Review

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	existReview, err := s.repo.FindReview(domain.RequestParams{
		Filter:  bson.M{"node_id": review.NodeID, "user_id": userIDPrimitive},
		Options: domain.Options{Limit: 1},
	})
	if err != nil {
		return nil, err
	}

	if len(existReview.Data) > 0 {
		updateReview := &model.ReviewInput{
			Rate:   review.Rate,
			Review: review.Review,
		}
		result, err = s.UpdateReview(existReview.Data[0].ID.Hex(), userID, updateReview)
	} else {
		result, err = s.repo.CreateReview(userID, review)

		// set user stat
		if err == nil {
			_, _ = s.userService.SetStat(userID, model.UserStat{AddReview: 1})
		}
	}

	return result, err
}

func (s *ReviewService) UpdateReview(id string, userID string, review *model.ReviewInput) (*model.Review, error) {
	return s.repo.UpdateReview(id, userID, review)
}

func (s *ReviewService) DeleteReview(id string) (*model.Review, error) {
	result, err := s.repo.DeleteReview(id)

	// set user stat
	if err == nil {
		_, _ = s.userService.SetStat(result.UserID.Hex(), model.UserStat{AddReview: -1})
	}

	return result, err
}
