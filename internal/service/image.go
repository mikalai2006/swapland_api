package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageService struct {
	repo        repository.Image
	imageConfig config.IImageConfig
}

func NewImageService(repo repository.Image, imageConfig config.IImageConfig) *ImageService {
	return &ImageService{repo: repo, imageConfig: imageConfig}
}

func (s *ImageService) FindImage(params domain.RequestParams) (domain.Response[model.Image], error) {
	return s.repo.FindImage(params)
}

func (s *ImageService) GetImage(id string) (model.Image, error) {
	return s.repo.GetImage(id)
}

func (s *ImageService) GetImageDirs(id string) ([]interface{}, error) {
	return s.repo.GetImageDirs(id)
}
func (s *ImageService) CreateImage(userID string, image *model.ImageInput) (model.Image, error) {
	var result model.Image

	if image.Service == "user" {
		userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return result, err
		}

		existImage, err := s.repo.FindImage(domain.RequestParams{Filter: bson.D{
			{"user_id", userIDPrimitive},
			{"service", image.Service},
			{"service_id", image.ServiceID},
		}})
		if err != nil {
			return result, err
		}

		if len(existImage.Data) > 0 {
			for i, _ := range existImage.Data {
				_, _ = s.DeleteImage(existImage.Data[i].ID.Hex())
				// if err != nil {
				// 	return result, err
				// }
			}
		}

	}
	result, err := s.repo.CreateImage(userID, image)

	return result, err
}

func (s *ImageService) DeleteImage(id string) (model.Image, error) {
	result := model.Image{}
	imageForRemove, err := s.GetImage(id)
	if err != nil {
		return result, err
	}
	if imageForRemove.Service == "" {
		return result, errors.New("not found item for remove")
	} else {
		pathOfRemove := fmt.Sprintf("public/%s/%s", imageForRemove.UserID.Hex(), imageForRemove.Service)

		if imageForRemove.ServiceID != "" {
			pathOfRemove = fmt.Sprintf("%s/%s", pathOfRemove, imageForRemove.ServiceID)
		}

		pathRemove := fmt.Sprintf("%s/%s%s", pathOfRemove, imageForRemove.Path, imageForRemove.Ext)
		os.Remove(pathRemove)
		// if err != nil {
		// 	return result, err
		// }

		// // remove srcset.
		// for i := range s.imageConfig.Sizes {
		// 	dataImg := s.imageConfig.Sizes[i]
		// 	pathRemove = fmt.Sprintf("%s/%v-%s%s", pathOfRemove, dataImg.Size, imageForRemove.Path, imageForRemove.Ext) // ".webp"
		// 	// fmt.Println("pathRemove2=", pathRemove)
		// 	os.Remove(pathRemove)
		// 	// if err != nil {
		// 	// 	return result, err
		// 	// }
		// }

		isEmpty, err := utils.IsEmptyDir(pathOfRemove)
		if err != nil {
			return result, err
		}
		if isEmpty {
			err = os.Remove(pathOfRemove)
			if err != nil {
				return result, err
			}
		}

		// pathRemove = fmt.Sprintf("%s/xs-%s", pathOfRemove, imageForRemove.Path)
		// err = os.Remove(pathRemove)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// }
		// pathRemove = fmt.Sprintf("%s/lg-%s", pathOfRemove, imageForRemove.Path)
		// err = os.Remove(pathRemove)
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// }
	}

	return s.repo.DeleteImage(id)
}
