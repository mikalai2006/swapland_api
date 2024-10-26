package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/pkg/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *HandlerV1) registerAuth(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.POST("/sign-up", h.SignUp)
	auth.POST("/sign-in", h.SignIn)
	auth.POST("/logout", h.Logout)
	auth.POST("/refresh", h.tokenRefresh)
	auth.GET("/refresh", h.tokenRefresh)
	auth.GET("/verification/:code", h.SetUserFromRequest, h.VerificationAuth)
	auth.GET("/iam", h.SetUserFromRequest, h.getIam)
}

func (h *HandlerV1) getIam(c *gin.Context) {
	appG := app.Gin{C: c}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	// TODO get token from body data.
	// var input *domain.RefreshInput

	// if err := c.BindJSON(&input); err != nil {
	// 	appG.Response(http.StatusBadRequest, err, nil)
	// 	return
	// }
	// fmt.Println("ID=", userID)

	users, err := h.Services.User.Iam(userID)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// get auth data for user
	authData, err := h.Services.GetAuth(users.UserID.Hex())
	if err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if !authData.ID.IsZero() {
		users.Md = authData.MaxDistance
		users.Roles = authData.Roles
		// fmt.Println("authData", authData)
	}

	// // implementation max distance.
	// md, err := middleware.GetMaxDistance(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// users.Md = md

	// // implementation roles for user.
	// roles, err := middleware.GetRoles(c)
	// if err != nil {
	// 	appG.ResponseError(http.StatusUnauthorized, err, nil)
	// 	return
	// }
	// users.Roles = roles

	c.JSON(http.StatusOK, users)
}

// @Summary SignUp
// @Tags auth
// @Description Create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body domain.Auth true "account info"
// @Success 200 {integer} 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-up [post].
func (h *HandlerV1) SignUp(c *gin.Context) {
	appG := app.Gin{C: c}

	lang := c.Query("lang")
	if lang == "" {
		lang = h.i18n.Default
	}

	var input *domain.SignInInput
	if err := c.BindJSON(&input); err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	input.Strategy = "local"

	// Check exist auth
	existAuth, err := h.Services.Authorization.ExistAuth(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}
	if !existAuth.ID.IsZero() {
		appG.ResponseError(http.StatusBadRequest, errors.New("exist account"), nil)
		return
	}

	id, err := h.Services.Authorization.CreateAuth(input)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	primitiveID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	// create default
	// avatar := fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon", id)

	newUser := model.User{
		// Avatar: avatar,
		UserID: primitiveID,
		Login:  input.Login,
		Name:   input.Login,
		Roles:  []string{"user"},
		Lang:   lang,
		Md:     50,
	}
	document, err := h.Services.User.CreateUser(id, &newUser)
	if err != nil {
		appG.ResponseError(http.StatusInternalServerError, err, nil)
		return
	}

	c.JSON(http.StatusOK, document)
}

// @Summary SignIn
// @Tags auth
// @Description Login user
// @ID signin-account
// @Accept json
// @Produce json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {integer} 1
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /auth/sign-in [post].
func (h *HandlerV1) SignIn(c *gin.Context) {
	appG := app.Gin{C: c}
	// jwt_cookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	// fmt.Println("+++++++++++++")
	// fmt.Printf("%s = %s",h.auth.NameCookieRefresh, jwt_cookie)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input *domain.SignInInput

	if err := c.BindJSON(&input); err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if input.Strategy == "" {
		input.Strategy = "local"
	}

	if input.Email == "" && input.Login == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("request must be with email or login"), nil)
		return
	}
	// if input.Password == "" {
	// 	appG.ResponseError(http.StatusBadRequest, errors.New("empty password"), nil)
	// 	return
	// }

	if input.Strategy == "local" {
		tokens, err := h.Services.Authorization.SignIn(input)
		if err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
		c.SetCookie(h.auth.NameCookieRefresh, tokens.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)

		c.JSON(http.StatusOK, domain.ResponseTokens{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresIn:    tokens.ExpiresIn,
		})
	}
	// else {
	// 	fmt.Print("JWT auth")
	// }
	// session.Set(userkey, input.Username)
	// session.Save()
}

// @Summary User Refresh Tokens
// @Tags users-auth
// @Description user refresh tokens
// @Accept  json
// @Produce  json
// @Param input body domain.RefreshInput true "sign up info"
// @Success 200 {object} domain.ResponseTokens
// @Failure 400,404 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Failure default {object} domain.ErrorResponse
// @Router /users/auth/refresh [post].
func (h *HandlerV1) tokenRefresh(c *gin.Context) {
	appG := app.Gin{C: c}
	jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	// fmt.Sprintf("refresh Cookie %s = %s", h.auth.NameCookieRefresh, jwtCookie)
	// cookie_header := c.GetHeader("cookie")
	// fmt.Println("refresh Cookie_header = ", cookie_header)
	// fmt.Println("+++++++++++++")
	// session := sessions.Default(c)
	var input domain.RefreshInput

	// if jwtCookie == "" {
	if err := c.BindJSON(&input); err != nil {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	// } else {
	// 	input.Token = jwtCookie
	// }
	// fmt.Println("refresh input.Token  = ", input.Token)
	// fmt.Println("jwtCookie  = ", jwtCookie)
	if input.Token == "" {
		input.Token = jwtCookie
	}

	if input.Token == "" && jwtCookie == "" {
		appG.ResponseError(http.StatusUnauthorized, errors.New("not found token"), nil)
		// c.JSON(http.StatusOK, gin.H{})
		// c.AbortWithStatus(http.StatusOK)
		return
	}

	res, err := h.Services.Authorization.RefreshTokens(input.Token)
	if err != nil && err != mongo.ErrNoDocuments {
		appG.ResponseError(http.StatusUnauthorized, err, nil)
		return
	}
	if err == mongo.ErrNoDocuments {
		c.SetCookie(h.auth.NameCookieRefresh, "", -1, "/", c.Request.URL.Hostname(), false, true)
	} else {
		c.SetCookie(h.auth.NameCookieRefresh, res.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)
	}

	// userData, err := h.services.User.FindUser(domain.RequestParams{Filter: bson.D{{"user_id": res.}}})
	// if err != nil {
	// 	appG.ResponseError(http.StatusBadRequest, err, nil)
	// 	return
	// }

	// c.SetCookie(h.auth.NameCookieRefresh, res.RefreshToken, int(h.auth.RefreshTokenTTL.Seconds()), "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, domain.ResponseTokens{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
		ExpiresInR:   res.ExpiresInR,
	})
}

func (h *HandlerV1) Logout(c *gin.Context) {
	// session := sessions.Default(c)
	// session.Delete(userkey)
	// if err := session.Save(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	// 	return
	// }
	appG := app.Gin{C: c}

	var input domain.RefreshInput

	jwtCookie, _ := c.Cookie(h.auth.NameCookieRefresh)
	if jwtCookie == "" {
		if err := c.BindJSON(&input); err != nil {
			appG.ResponseError(http.StatusBadRequest, err, nil)
			return
		}
	} else {
		input.Token = jwtCookie
	}

	if input.Token == "" && jwtCookie == "" {
		c.JSON(http.StatusOK, gin.H{})
		c.AbortWithStatus(http.StatusOK)
		return
	}

	_, err := h.Services.Authorization.RemoveRefreshTokens(input.Token)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	c.SetCookie(h.auth.NameCookieRefresh, "", -1, "/", c.Request.URL.Hostname(), false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (h *HandlerV1) VerificationAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	code := c.Param("code")
	if code == "" {
		appG.ResponseError(http.StatusBadRequest, errors.New("code empty"), nil)
		return
	}

	userID, err := middleware.GetUID(c)
	if err != nil {
		appG.ResponseError(http.StatusBadRequest, err, nil)
		return
	}

	if er := h.Services.Authorization.VerificationCode(userID, code); er != nil {
		appG.ResponseError(http.StatusBadRequest, er, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
