package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"github.com/mikalai2006/swapland-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *HandlerV1) registerAddress(router *gin.RouterGroup) {
	address := router.Group("/address")
	address.GET("/", h.FindAddress)
	address.POST("/", h.CreateAddress)
	address.POST("/list/", h.CreateListAddress)
}

func (h *HandlerV1) CreateAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	lang := c.Query("lang")
	if lang == "" {
		lang = h.i18n.Default
	}
	fmt.Println("lang", lang)

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input *domain.AddressInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	if input.Lang == "" {
		input.Lang = lang
	}

	existOsmID, err := h.Services.Address.FindAddress(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"osm_id", input.OsmID}},
	})
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if len(existOsmID.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("items exists"), nil)
		return
	}

	address, err := h.Services.Address.CreateAddress(userID, input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, address)
}

func (h *HandlerV1) CreateListAddress(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var input []*domain.AddressInput
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if len(input) == 0 {
		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
		return
	}

	var result []*domain.Address
	for i := range input {
		// existOsmID, err := h.services.Address.FindAddress(domain.RequestParams{
		// 	Options: domain.Options{Limit: 1},
		// 	Filter:  bson.D{{"osm_id", input[i].OsmID}},
		// })
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }

		// if len(existOsmID.Data) == 0 {
		address, err := h.Services.Address.CreateAddress(userID, input[i])
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, address)
		// }

	}

	c.JSON(http.StatusOK, result)
}

// @Summary Address Get all Address
// @Security ApiKeyAuth
// @Tags address
// @Description get all Address
// @ModuleID address
// @Accept  json
// @Produce  json
// @Success 200 {object} []domain.Address
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/address [get].
func (h *HandlerV1) GetAllAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Address{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	addresses, err := h.Services.Address.GetAllAddress(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, addresses)
}

// @Summary Address by params
// @Security ApiKeyAuth
// @Tags address
// @Description Input params for search Addresses
// @ModuleID address
// @Accept  json
// @Produce  json
// @Param input query Address true "params for search Address"
// @Success 200 {object} []domain.Address
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/address [get].
func (h *HandlerV1) FindAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, domain.Address{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	addresses, err := h.Services.Address.FindAddress(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (h *HandlerV1) GetAddressByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateAddress(c *gin.Context) {

}

func (h *HandlerV1) DeleteAddress(c *gin.Context) {

}
