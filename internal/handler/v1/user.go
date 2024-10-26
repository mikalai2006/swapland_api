package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"github.com/mikalai2006/swapland-api/pkg/app"
)

func (h *HandlerV1) RegisterUser(router *gin.RouterGroup) {
	user := router.Group("/user")
	user.POST("", h.CreateUser)
	user.POST("/find", h.FindUser)
	user.GET("/:id", h.GetUser)
	user.DELETE("/:id", h.DeleteUser)
	user.PATCH("/:id", h.UpdateUser)
	user.PATCH("/:id/location", h.UpdateLocation)
}

// @Summary Get user by Id
// @Tags user
// @Description get user info
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [get].
func (h *HandlerV1) GetUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	user, err := h.Services.User.GetUser(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// // get auth data for user
	// authData, err := h.services.GetAuth(user.UserID.Hex())
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// if !authData.ID.IsZero() {
	// 	user.Md = authData.MaxDistance
	// 	user.Roles = authData.Roles
	// }

	c.JSON(http.StatusOK, user)
}

// type InputUser struct {
// 	domain.RequestParams
// 	model.User
// }

// @Summary Find few users
// @Security ApiKeyAuth
// @Tags user
// @Description Input params for search users
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param input query model.UserInput true "params for search users"
// @Success 200 {object} []model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user [get].
func (h *HandlerV1) FindUser(c *gin.Context) {
	appG := app.Gin{C: c}

	params, err := utils.GetParamsFromRequest(c, model.UserInput{}, &h.i18n)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	users, err := h.Services.User.FindUser(params)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *HandlerV1) CreateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var input *model.User
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	user, er := h.Services.User.CreateUser(userID, input)
	if er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Security ApiKeyAuth
// @Tags user
// @Description Delete user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} []model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [delete].
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")

	// var input model.User
	// if err := c.BindJSON(&input); err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())

	// 	return
	// }

	user, err := h.Services.User.DeleteUser(id) // , input
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update user
// @Security ApiKeyAuth
// @Tags user
// @Description Update user
// @ModuleID user
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Param input body model.User true "body for update user"
// @Success 200 {object} []model.User
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /api/user/{id} [put].
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}

	id := c.Param("id")
	fmt.Println("UpdateUser: ", id)

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	var input model.User
	if er := c.Bind(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	fmt.Println("input: ", input)
	// fmt.Println("UpdateUser input: ", input)
	if input.Location.Lat != 0 && input.Location.Lon != 0 {
		// Get address.
		pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", input.Location.Lat, input.Location.Lon))
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
		r.Header.Add("User-Agent", "a127.0.0.1")

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		defer resp.Body.Close()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		var bodyResponse domain.ResponseNominatim
		if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
			appG.ResponseError(http.StatusBadRequest, e, nil)
			return
		}

		if bodyResponse.OsmID != 0 {
			// Check address in to db
			input.Location.OsmID = fmt.Sprintf("%s/%d", bodyResponse.OsmType, bodyResponse.OsmID)
			input.Location.Address = bodyResponse.Address

			// if bodyResponse.Name == "" {
			// 	arrStr := strings.Split(bodyResponse.Address, ",")
			// 	nameNode := ""
			// 	if len(arrStr) >= 2 {
			// 		nameNode = fmt.Sprintf("%s, %s", arrStr[1], arrStr[0])
			// 	} else {
			// 		nameNode = arrStr[0]
			// 	}
			// 	input.Name = strings.TrimSpace(nameNode)
			// } else {
			// 	input.Name = bodyResponse.Name
			// }
		}
		//  else {
		// 	fmt.Println("not found osm", bodyResponse.OsmID)
		// }
	}

	user, err := h.Services.User.UpdateUser(id, &input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var imageInput = &model.ImageInput{}
	imageInput.Service = "user"
	imageInput.ServiceID = user.ID.Hex()
	imageInput.UserID = userID
	imageInput.Dir = "user"

	// fmt.Println("imageInput: ", imageInput)

	paths, err := utils.UploadResizeMultipleFile(c, imageInput, "images", &h.imageConfig)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	var result []model.Image
	for i := range paths {
		imageInput.Path = paths[i].Path
		imageInput.Ext = paths[i].Ext
		image, err := h.Services.Image.CreateImage(userID, imageInput)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		result = append(result, image)
	}

	c.JSON(http.StatusOK, user)
}

func (h *HandlerV1) UpdateLocation(c *gin.Context) {
	appG := app.Gin{C: c}

	// id := c.Param("id")
	// fmt.Println("UpdateLocation: ", id)

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}

	var inputUser model.User
	var input model.GeoLocation
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	fmt.Println("input: ", input)
	// fmt.Println("UpdateUser input: ", input)
	if input.Lat != 0 && input.Lon != 0 {
		// Get address.
		pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", input.Lat, input.Lon))
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
		r.Header.Add("User-Agent", "a127.0.0.1")

		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		defer resp.Body.Close()

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		var bodyResponse domain.ResponseNominatim
		if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
			appG.ResponseError(http.StatusBadRequest, e, nil)
			return
		}

		if bodyResponse.OsmID != 0 {
			// Check address in to db
			inputUser.Location.Lat = input.Lat
			inputUser.Location.Lon = input.Lon
			inputUser.Location.OsmID = fmt.Sprintf("%s/%d", bodyResponse.OsmType, bodyResponse.OsmID)
			inputUser.Location.Address = bodyResponse.Address

			// if bodyResponse.Name == "" {
			// 	arrStr := strings.Split(bodyResponse.Address, ",")
			// 	nameNode := ""
			// 	if len(arrStr) >= 2 {
			// 		nameNode = fmt.Sprintf("%s, %s", arrStr[1], arrStr[0])
			// 	} else {
			// 		nameNode = arrStr[0]
			// 	}
			// 	input.Name = strings.TrimSpace(nameNode)
			// } else {
			// 	input.Name = bodyResponse.Name
			// }
		}
		//  else {
		// 	fmt.Println("not found osm", bodyResponse.OsmID)
		// }
	}

	user, err := h.Services.User.UpdateUser(userID, &inputUser)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, user)
}
