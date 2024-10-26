package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"github.com/mikalai2006/swapland-api/pkg/app"
)

func (h *HandlerV1) registerTrack(router *gin.RouterGroup) {
	track := router.Group("/track")
	track.GET("/", h.FindTrack)
	track.POST("/", h.CreateTrack)
	track.POST("/list/", h.CreateListTrack)
}

func (h *HandlerV1) CreateTrack(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *domain.Track
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	track, err := h.Services.Track.CreateTrack(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, track)
}

func (h *HandlerV1) CreateListTrack(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.Track
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Track
	for i := range input {
		track, err := h.Services.Track.CreateTrack(userID, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, track)

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Track Get all Tracks
// @Security ApiKeyAuth
// @Tags Track
// @Description get all Tracks
// @ModuleID Track
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Track
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Track [get].
func (h *HandlerV1) GetAllTrack(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Track{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tracks, err := h.Services.Track.GetAllTrack(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tracks)
}

// @Summary Find Tracks by params
// @Security ApiKeyAuth
// @Tags Track
// @Description Input params for search Tracks
// @ModuleID Track
// @Accept  json
// @Produce  json
// @Param input query TrackInput true "params for search Track"
// @Success 200 {object} []domain.Track
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Track [get].
func (h *HandlerV1) FindTrack(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.TrackInputData{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Tracks, err := h.Services.Track.FindTrack(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Tracks)
}

func (h *HandlerV1) GetTrackByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateTrack(c *gin.Context) {

}

func (h *HandlerV1) DeleteTrack(c *gin.Context) {

}
