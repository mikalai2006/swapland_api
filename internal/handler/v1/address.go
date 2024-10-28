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
	"github.com/mikalai2006/swapland-api/pkg/app"
)

func (h *HandlerV1) registerAddress(router *gin.RouterGroup) {
	address := router.Group("/address")
	address.POST("/", h.CreateAddress)
	address.POST("/find", h.FindAddress)
	address.POST("/address", h.GetAddress)
	address.PATCH("/:id", h.SetUserFromRequest, h.UpdateAddress)
}

func (h *HandlerV1) CreateAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	lang := c.Query("lang")
	if lang == "" {
		lang = h.i18n.Default
	}

	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}

	// var input *model.Address
	// if er := c.BindJSON(&input); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// if input.Lang == "" {
	// 	input.Lang = lang
	// }
	var address *model.Address
	var input *model.GeoCoordinates
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	var bodyResponse domain.ResponseNominatim

	fmt.Println("latlon: ", *input.Lat, *input.Lon)

	if input.Lat != nil && input.Lon != nil {
		// Get address.
		bodyResponse, err = GetAddress(c, input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		// existOsmID, err := h.Services.Address.FindAddress(&model.AddressFilter{Lat: *input.Lat, Lon: *input.Lon})
		// if err != nil {
		// 	appG.ResponseError(http.StatusBadRequest, err, nil)
		// 	return
		// }

		// if len(existOsmID.Data) > 0 {
		// 	// appG.ResponseError(http.StatusBadRequest, errors.New("items exists"), nil)
		// 	// return
		// 	address = &existOsmID.Data[0]
		// } else {
		address, err = h.Services.Address.CreateAddress(userID, bodyResponse)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		// }
	}
	c.JSON(http.StatusOK, address)
}

func GetAddress(c *gin.Context, input *model.GeoCoordinates) (domain.ResponseNominatim, error) {
	var response domain.ResponseNominatim

	pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=ru", *input.Lat, *input.Lon))
	if err != nil {
		// appG.ResponseError(http.StatusBadRequest, err, nil)
		return response, err
	}
	r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody)
	r.Header.Add("User-Agent", "a127.0.0.1")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		// appG.ResponseError(http.StatusBadRequest, err, nil)
		return response, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// appG.ResponseError(http.StatusBadRequest, err, nil)
		return response, err
	}

	if e := json.Unmarshal(bytes, &response); e != nil {
		// appG.ResponseError(http.StatusBadRequest, e, nil)
		return response, err
	}

	return response, err
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

	var input *model.AddressFilter
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	addresses, err := h.Services.Address.FindAddress(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func (h *HandlerV1) GetAddress(c *gin.Context) {
	appG := app.Gin{C: c}

	var input *model.GeoCoordinates
	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}
	var bodyResponse domain.ResponseNominatim

	if input.Lat != nil && input.Lon != nil {
		// Get address.
		pathRequest, err := url.Parse(fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%f&lon=%f&format=json&accept-language=none", *input.Lat, *input.Lon))
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

		if e := json.Unmarshal(bytes, &bodyResponse); e != nil {
			appG.ResponseError(http.StatusBadRequest, e, nil)
			return
		}
	}

	c.JSON(http.StatusOK, bodyResponse)
}

func (h *HandlerV1) UpdateAddress(c *gin.Context) {
	appG := app.Gin{C: c}
	userID, err := middleware.GetUID(c)
	if err != nil {
		// c.AbortWithError(http.StatusUnauthorized, err)
		appG.ResponseError(http.StatusUnauthorized, err, gin.H{"hello": "world"})
		return
	}
	id := c.Param("id")

	// var a map[string]json.RawMessage //  map[string]interface{}
	// if er := c.ShouldBindBodyWith(&a, binding.JSON); er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }
	// data, er := utils.BindJSON2[model.Address](a)
	// if er != nil {
	// 	appG.ResponseError(http.StatusBadRequest, er, nil)
	// 	return
	// }

	var address *model.Address
	var input *model.GeoCoordinates

	if er := c.BindJSON(&input); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	if input.Lat != nil && input.Lon != nil {
		// Get address.
		bodyResponse, err := GetAddress(c, input)

		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		address, err = h.Services.Address.UpdateAddress(id, userID, bodyResponse)
		if err != nil {
			appG.ResponseError(http.StatusInternalServerError, err, nil)
			return
		}

		// update products with address.
		products, err := h.Services.Product.FindProduct(&model.ProductFilter{AddressId: &id})
		if err != nil {
			appG.ResponseError(http.StatusInternalServerError, err, nil)
			return
		}

		for i := range products.Data {
			_, err = h.Services.Product.UpdateProduct(products.Data[i].ID.Hex(), userID,
				&model.Product{
					Lon: address.Lon,
					Lat: address.Lat,
				})
			if err != nil {
				appG.ResponseError(http.StatusInternalServerError, err, nil)
				return
			}
		}
		// fmt.Println("products: ", len(products.Data))
	}

	c.JSON(http.StatusOK, address)
}

func (h *HandlerV1) DeleteAddress(c *gin.Context) {

}
