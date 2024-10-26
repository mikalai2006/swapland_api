package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"github.com/mikalai2006/swapland-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *HandlerV1) registerSubscribe(router *gin.RouterGroup) {
	subscribe := router.Group("/subscribe")
	subscribe.GET("", h.FindSubscribe)
	subscribe.POST("", h.CreateSubscribe)
	subscribe.PATCH("/:id", h.UpdateSubscribe)
	subscribe.DELETE("/:id", h.DeleteSubscribe)
}

func (h *HandlerV1) CreateSubscribe(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	var data *model.SubscribeInput
	if er := c.BindJSON(&data); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// Check exists Subscribefor node and user
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
	}
	// subUserIDPrimitive, err := primitive.ObjectIDFromHex(data.SubUserID)
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// }
	SubscribeExist, err := h.Services.Subscribe.FindSubscribe(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  bson.D{{"sub_user_id", data.SubUserID}, {"user_id", userIDPrimitive}},
	},
	)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if len(SubscribeExist.Data) > 0 {
		appG.ResponseError(http.StatusBadRequest, model.ErrSubscribeExist, nil)
		return
	}

	Subscribe, err := h.Services.Subscribe.CreateSubscribe(userID, data)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Subscribe)
}

// @Summary Find Subscribes by params
// @Security ApiKeyAuth
// @Subscribes Subscribe
// @Description Input params for search Subscribes
// @ModuleID Subscribe
// @Accept  json
// @Produce  json
// @Param input query SubscribeInput true "params for search Subscribe"
// @Success 200 {object} []model.Subscribe
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/Subscribe [get].
func (h *HandlerV1) FindSubscribe(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.SubscribeInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	Subscribes, err := h.Services.Subscribe.FindSubscribe(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, Subscribes)
}

func (h *HandlerV1) GetSubscribeByID(c *gin.Context) {

}

func (h *HandlerV1) UpdateSubscribe(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	var data *model.Subscribe
	if er := c.BindJSON(&data); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	// fmt.Println(data)

	document, err := h.Services.Subscribe.UpdateSubscribe(id, userID, data)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *HandlerV1) DeleteSubscribe(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	if id == "" {
		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
		return
	}

	user, err := h.Services.Subscribe.DeleteSubscribe(id) // , input
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
