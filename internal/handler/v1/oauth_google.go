package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
}

func (h *HandlerV1) registerGoogleOAuth(router *gin.RouterGroup) {
	router.GET("/google", h.OAuthGoogle)
	router.GET("/google/me", h.MeGoogle)
}

func (h *HandlerV1) OAuthGoogle(c *gin.Context) {
	appG := app.Gin{C: c}

	urlReferer := c.Request.Referer()
	scope := strings.Join(h.oauth.GoogleScopes, " ")

	pathRequest, err := url.Parse(h.oauth.GoogleAuthURI)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	parameters := url.Values{}
	parameters.Add("client_id", h.oauth.GoogleClientID)
	parameters.Add("redirect_uri", h.oauth.GoogleRedirectURI)
	parameters.Add("scope", scope)
	parameters.Add("response_type", "code")
	parameters.Add("prompt", "select_account")
	parameters.Add("state", urlReferer)
	fmt.Println("client URL(state): ", urlReferer)
	// if urlReferer == "" {
	// 	urlReferer = "file:///android_asset/auth.html"
	// } else {
	// 	urlReferer = fmt.Sprintf("%s%s", urlReferer, "app")
	// }
	// fmt.Println("urlReferer= ", urlReferer)

	pathRequest.RawQuery = parameters.Encode()
	// fmt.Println("Google auth1::: ", pathRequest.String())
	c.Redirect(http.StatusFound, pathRequest.String())
}

func (h *HandlerV1) MeGoogle(c *gin.Context) {
	appG := app.Gin{C: c}

	code := c.Query("code")
	clientURL := c.Query("state")

	if code == "" {
		pathRequest, err := url.Parse(clientURL)
		if err != nil {
			// c.AbortWithError(http.StatusBadRequest, err)
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		parameters := url.Values{}
		errorClient := c.Query("error")
		parameters.Add("error", errorClient)
		pathRequest.RawQuery = parameters.Encode()
		c.Redirect(http.StatusFound, pathRequest.String())
		// // c.AbortWithError(http.StatusBadRequest, errors.New("no correct code"))
		// appG.ResponseError(http.StatusBadRequest, errors.New("no correct code"), nil)
		return
	}
	// fmt.Println("error!!!!!!!!!!!!!!!!!!!!!!!!!!")

	pathRequest, err := url.Parse(h.oauth.GoogleTokenURI)
	if err != nil {
		panic("boom")
	}
	parameters := url.Values{}
	parameters.Set("client_id", h.oauth.GoogleClientID)
	parameters.Set("redirect_uri", h.oauth.GoogleRedirectURI)
	parameters.Set("client_secret", h.oauth.GoogleClientSecret)
	parameters.Set("code", code)
	parameters.Set("grant_type", "authorization_code")

	req, _ := http.NewRequestWithContext(c, http.MethodPost, pathRequest.String(), strings.NewReader(parameters.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	defer resp.Body.Close()

	token := struct {
		AccessToken string `json:"access_token"`
	}{}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if er := json.Unmarshal(bytes, &token); er != nil { // Parse []byte to go struct pointer
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	pathRequest, err = url.Parse(h.oauth.GoogleUserinfoURI)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	r, _ := http.NewRequestWithContext(c, http.MethodGet, pathRequest.String(), http.NoBody) // URL-encoded payload
	bearerToken := fmt.Sprintf("Bearer %s", token.AccessToken)
	r.Header.Add("Authorization", bearerToken)
	// r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = http.DefaultClient.Do(r)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	defer resp.Body.Close()

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	var bodyResponse GoogleUserInfo
	if e := json.Unmarshal(bytes, &bodyResponse); e != nil { // Parse []byte to go struct pointer
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, e, nil)
		return
	}

	input := &domain.SignInInput{
		Login:       bodyResponse.Email,
		Strategy:    "jwt",
		Password:    "",
		GoogleID:    bodyResponse.Sub,
		MaxDistance: 100,
	}
	fmt.Println("Google auth2::: ", input)

	user, err := h.Services.Authorization.ExistAuth(input)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// Create new user, if not exist
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
	// if user.Login == "" {
	// 	a, err := h.services.Authorization.CreateAuth(input)
	// 	fmt.Println("Google auth3::: ", a)
	// 	if err != nil {
	// 		// c.AbortWithError(http.StatusBadRequest, err)
	// 		appG.ResponseError(http.StatusBadRequest, err, nil)
	// 		return
	// 	}
	// }

	tokens, err := h.Services.Authorization.SignIn(input)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	h.Services.User.UpdateUser(userID, &model.User{Online: true})

	// TODO.
	// if clientURL == "" {
	// 	clientURL = "https://www.poihub.ru/app"
	// } else {

	// }
	clientURL = "swapland://profile"
	fmt.Println("clientURL: ", clientURL)

	pathRequest, err = url.Parse(clientURL)
	if err != nil {
		// c.AbortWithError(http.StatusBadRequest, err)
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	parameters = url.Values{}
	parameters.Add("token", tokens.AccessToken)
	// if len(clientURL) == 0 || clientURL == "http://localhost:8081/" {
	// }
	parameters.Add("rt", tokens.RefreshToken)
	parameters.Add("exp", strconv.FormatInt(tokens.ExpiresIn, 10))
	parameters.Add("expr", strconv.FormatInt(tokens.ExpiresInR, 10))
	pathRequest.RawQuery = parameters.Encode()
	c.SetCookie(h.auth.NameCookieRefresh, tokens.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)
	c.Redirect(http.StatusFound, pathRequest.String())
}
