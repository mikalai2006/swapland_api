package v1

// import (
// 	"errors"
// 	"math"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/mikalai2006/swapland-api/graph/model"
// 	"github.com/mikalai2006/swapland-api/internal/domain"
// 	"github.com/mikalai2006/swapland-api/internal/middleware"
// 	"github.com/mikalai2006/swapland-api/internal/utils"
// 	"github.com/mikalai2006/swapland-api/pkg/app"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func (h *HandlerV1) registerNodeVote(router *gin.RouterGroup) {
// 	nodeVote := router.Group("/node_vote")
// 	nodeVote.GET("", h.FindNodeVote)
// 	nodeVote.POST("", h.CreateNodeVote)
// 	nodeVote.POST("/list", h.CreateNodeVoteList)
// 	nodeVote.PATCH("/:id", h.UpdateNodeVote)
// 	nodeVote.DELETE("/:id", h.DeleteNodeVote)
// }

// func (h *HandlerV1) CreateNodeVote(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	var input *model.NodeVote
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	NodeVote, err := h.CreateOrExistNodeVote(c, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, NodeVote)
// }

// func (h *HandlerV1) CreateNodeVoteList(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	var input []*model.NodeVote
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	if len(input) == 0 {
// 		appG.ResponseError(http.StatusBadRequest, errors.New("list must be with element(s)"), nil)
// 		return
// 	}

// 	var result []*model.NodeVote
// 	for i := range input {
// 		NodeVote, err := h.CreateOrExistNodeVote(c, input[i])
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return
// 		}

// 		result = append(result, NodeVote)
// 	}

// 	c.JSON(http.StatusOK, result)
// }

// func (h *HandlerV1) CreateOrExistNodeVote(c *gin.Context, input *model.NodeVote) (*model.NodeVote, error) {
// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return nil, err
// 	}
// 	var result *model.NodeVote

// 	// nodedataIDPrimitive, err := primitive.ObjectIDFromHex(input.NodedataID)
// 	// if err != nil {
// 	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
// 	// 	return result, err
// 	// }
// 	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}

// 	// check exist node
// 	existNodes, err := h.services.Node.FindNode(domain.RequestParams{
// 		Filter:  bson.D{{"_id", input.NodeID}},
// 		Options: domain.Options{Limit: 1},
// 	})
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}
// 	if len(existNodes.Data) == 0 {
// 		// appG.ResponseError(http.StatusBadRequest, errors.New("not found nodedata"), nil)
// 		return result, nil
// 	}
// 	input.NodeUserID = existNodes.Data[0].UserID

// 	// check exist vote
// 	existNodeVote, err := h.services.NodeVote.FindNodeVote(domain.RequestParams{
// 		Options: domain.Options{Limit: 1},
// 		Filter: bson.D{
// 			{"value", input.Value},
// 			{"node_id", input.NodeID},
// 			{"user_id", userIDPrimitive},
// 		},
// 	})
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}
// 	if len(existNodeVote.Data) > 0 {
// 		// appG.ResponseError(http.StatusBadRequest, model.ErrNodeVoteExistValue, nil)
// 		// return &existNodeVote.Data[0], nil
// 		result = &existNodeVote.Data[0]
// 	} else {
// 		result, err = h.services.NodeVote.CreateNodeVote(userID, input)
// 		if err != nil {
// 			appG.ResponseError(http.StatusBadRequest, err, nil)
// 			return result, err
// 		}
// 	}

// 	// update counter votes node.
// 	votes, err := h.services.NodeVote.FindNodeVote(domain.RequestParams{
// 		Filter: bson.D{{"node_id", result.NodeID}},
// 	})
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}

// 	like := 0
// 	dlike := 0
// 	for i, _ := range votes.Data {
// 		if votes.Data[i].Value > 0 {
// 			like += 1
// 		} else {
// 			dlike += 1
// 		}
// 	}
// 	status := 100
// 	if dlike > 5 {
// 		status = -1
// 	}
// 	_, err = h.services.Node.UpdateNode(result.NodeID.Hex(), userID, &model.Node{
// 		NodeLike: model.NodeLike{
// 			Like:   int64(like),
// 			Dlike:  int64(math.Abs(float64(dlike))),
// 			Status: int64(status),
// 		},
// 	})
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return result, err
// 	}

// 	return result, nil
// }

// // // @Summary NodeVote Get all NodeVotes
// // // @Security ApiKeyAuth
// // // @NodeVotes NodeVote
// // // @Description get all NodeVotes
// // // @ModuleID NodeVote
// // // @Accept  json
// // // @Produce  json
// // // @Success 200 {object} []model.NodeVote
// // // @Failure 400,404 {object} domain.ErrorResponse
// // // @Failure 500 {object} domain.ErrorResponse
// // // @Failure default {object} domain.ErrorResponse
// // // @Router /api/node_vote [get].
// // func (h *HandlerV1) GetAllNodeVote(c *gin.Context) {
// // 	appG := app.Gin{C: c}

// // 	params, err := utils.GetParamsFromRequest(c, model.NodeVote{}, &h.i18n)
// // 	if err != nil {
// // 		appG.ResponseError(http.StatusBadRequest, err, nil)
// // 		return
// // 	}

// // 	NodeVotes, err := h.services.NodeVote.GetAllNodeVote(params)
// // 	if err != nil {
// // 		appG.ResponseError(http.StatusBadRequest, err, nil)
// // 		return
// // 	}

// // 	c.JSON(http.StatusOK, NodeVotes)
// // }

// // @Summary Find NodeVotes by params
// // @Security ApiKeyAuth
// // @NodeVotes NodeVote
// // @Description Input params for search NodeVotes
// // @ModuleID NodeVote
// // @Accept  json
// // @Produce  json
// // @Param input query NodeVoteInput true "params for search NodeVote"
// // @Success 200 {object} []model.NodeVote
// // @Failure 400,404 {object} domain.ErrorResponse
// // @Failure 500 {object} domain.ErrorResponse
// // @Failure default {object} domain.ErrorResponse
// // @Router /api/node_vote [get].
// func (h *HandlerV1) FindNodeVote(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	params, err := utils.GetParamsFromRequest(c, model.NodeVote{}, &h.i18n)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	NodeVotes, err := h.services.NodeVote.FindNodeVote(params)
// 	if err != nil {
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, NodeVotes)
// }

// func (h *HandlerV1) UpdateNodeVote(c *gin.Context) {

// 	appG := app.Gin{C: c}
// 	userID, err := middleware.GetUID(c)
// 	if err != nil {
// 		// c.AbortWithError(http.StatusUnauthorized, err)
// 		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
// 		return
// 	}
// 	id := c.Param("id")

// 	var input *model.NodeVoteInput
// 	if er := c.BindJSON(&input); er != nil {
// 		appG.ResponseError(http.StatusBadRequest, er, nil)
// 		return
// 	}

// 	document, err := h.services.NodeVote.UpdateNodeVote(id, userID, input)
// 	if err != nil {
// 		appG.ResponseError(http.StatusInternalServerError, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, document)
// }

// func (h *HandlerV1) DeleteNodeVote(c *gin.Context) {
// 	appG := app.Gin{C: c}

// 	id := c.Param("id")
// 	if id == "" {
// 		// c.AbortWithError(http.StatusBadRequest, errors.New("for remove need id"))
// 		appG.ResponseError(http.StatusBadRequest, errors.New("for remove need id"), nil)
// 		return
// 	}

// 	user, err := h.services.NodeVote.DeleteNodeVote(id) // , input
// 	if err != nil {
// 		// c.AbortWithError(http.StatusBadRequest, err)
// 		appG.ResponseError(http.StatusBadRequest, err, nil)
// 		return
// 	}

// 	c.JSON(http.StatusOK, user)
// }
