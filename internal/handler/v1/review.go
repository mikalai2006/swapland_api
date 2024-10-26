package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"github.com/mikalai2006/swapland-api/pkg/app"
)

func (h *HandlerV1) registerReview(router *gin.RouterGroup) {
	review := router.Group("/review")
	review.GET("/", h.FindReview)
	review.POST("", h.CreateReview)
	review.POST("/list", h.CreateReviewList)
}

func (h *HandlerV1) CreateReview(c *gin.Context) {
	appG := app.Gin{C: c}
	// userID, err := middleware.GetUID(c)
	// if err != nil {
	// 	// c.AbortWithError(http.StatusUnauthorized, err)
	// 	appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
	// 	return
	// }

	var input *model.Review
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	review, err := h.CreateOrExistReview(c, input) //h.services.Review.CreateReview(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *HandlerV1) CreateReviewList(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil || userID == "" {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*model.Review
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*model.Review
	for i := range input {
		review, err := h.CreateOrExistReview(c, input[i]) //h.services.Review.CreateReview(userID, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, review)
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Review Get all reviews
// @Security ApiKeyAuth
// @Tags review
// @Description get all reviews
// @ModuleID review
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Review
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/review [get].
func (h *HandlerV1) GetAllReview(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.Review{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	reviews, err := h.Services.Review.GetAllReview(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// @Summary Find reviews by params
// @Security ApiKeyAuth
// @Tags review
// @Description Input params for search reviews
// @ModuleID review
// @Accept  json
// @Produce  json
// @Param input query ReviewInput true "params for search review"
// @Success 200 {object} []domain.Review
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/review [get].
func (h *HandlerV1) FindReview(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.ReviewInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	reviews, err := h.Services.Review.FindReview(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *HandlerV1) GetReviewByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateReview(c *gin.Context) {

}

func (h *HandlerV1) DeleteReview(c *gin.Context) {

}

func (h *HandlerV1) CreateOrExistReview(c *gin.Context, input *model.Review) (*model.Review, error) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return nil, err
	}
	var result *model.Review

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }

	// existReviews, err := h.services.Review.FindReview(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1},
	// 	Filter:  bson.D{{"node_id", input.NodeID}, {"user_id", userIDPrimitive}},
	// })
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return result, err
	// }
	// if len(existReviews.Data) > 0 {
	// 	fmt.Println("existReviews =")
	// 	// appG.ResponseError(http.StatusBadRequest, model.ErrNodedataVoteExistValue, nil)
	// 	return &existReviews.Data[0], nil
	// }

	result, err = h.Services.Review.CreateReview(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return result, err
	}
	return result, nil
}
