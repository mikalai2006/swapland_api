package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VKBodyResponse struct {
	Response []struct {
		ID              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
	} `json:"response"`
}

func (h *HandlerV1) registerVkOAuth(router *gin.RouterGroup) {
	router.GET("/vk", h.OAuthVK)
	router.GET("/vk/me", h.MeVk)
}

func (h *HandlerV1) OAuthVK(c *gin.Context) {
	appG := app.Gin{C: c}

	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.VkScopes, "+")

	pathRequest, err := url.Parse(h.oauth.VkAuthURI)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", h.oauth.VkClientID)
	parameters.Add("redirect_uri", h.oauth.VkRedirectURI)
	parameters.Add("scope", scope)
	parameters.Add("response_type", "code")
	parameters.Add("state", urlReferer)

	pathRequest.RawQuery = parameters.Encode()
	c.Redirect(http.StatusFound, pathRequest.String())
}

func (h *HandlerV1) MeVk(c *gin.Context) {
	appG := app.Gin{C: c}

	code := c.Query("code")
	clientURL := c.Query("state")
	if code == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("no correct code"), nil)
		return
	}

	pathRequest, err := url.Parse(h.oauth.VkTokenURI)
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Set("client_id", h.oauth.VkClientID)
	parameters.Set("client_secret", h.oauth.VkClientSecret)
	parameters.Set("redirect_uri", h.oauth.VkRedirectURI)
	parameters.Set("code", code)

	req, _ := http.NewRequestWithContext(c, http.MethodPost, pathRequest.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	defer resp.Body.Close()

	token := struct {
		AccessToken string `json:"access_token"`
	}{}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if er := json.Unmarshal(bytes, &token); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	pathRequest, err = url.Parse(h.oauth.VkUserinfoURI)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	parameters = url.Values{}
	parameters.Set("access_token", token.AccessToken)
	parameters.Set("v", "5.131")

	req, _ = http.NewRequestWithContext(c, http.MethodPost, pathRequest.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var bodyResponse VKBodyResponse
	if er := json.Unmarshal(bytes, &bodyResponse); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	// fmt.Println("input Vk")
	input := &domain.SignInInput{
		Login:       bodyResponse.Response[0].FirstName,
		Strategy:    "jwt",
		Password:    "1",
		VkID:        fmt.Sprintf("%d", bodyResponse.Response[0].ID),
		MaxDistance: 100,
	}

	// fmt.Println("input Vk", input)
	user, err := h.Services.Authorization.ExistAuth(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var userID string
	if user.Login == "" {
		userID, err = h.Services.Authorization.CreateAuth(input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}

		primitiveID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		newUser := model.User{
			UserID: primitiveID,
			Login:  input.Login,
			Name:   input.Login,
			Roles:  []string{"user"},
			Lang:   h.i18n.Default,
		}
		_, err = h.Services.User.CreateUser(userID, &newUser)
		if err != nil {
			appG.ResponseError(http.StatusInternalServerError, err, nil)
			return
		}
	}

	tokens, err := h.Services.Authorization.SignIn(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	pathRequest, err = url.Parse(clientURL)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	parameters = url.Values{}
	parameters.Add("token", tokens.AccessToken)
	pathRequest.RawQuery = parameters.Encode()
	// c.Redirect(http.StatusMovedPermanently, path)
	c.SetCookie(h.auth.NameCookieRefresh, tokens.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, pathRequest.String())
}
