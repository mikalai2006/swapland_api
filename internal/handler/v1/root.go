package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"github.com/mikalai2006/swapland-api/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerV1 struct {
	db           *mongo.Database
	repositories *repository.Repositories
	Services     *service.Services
	oauth        config.OauthConfig
	auth         config.AuthConfig
	i18n         config.I18nConfig
	imageConfig  config.IImageConfig
	hub          service.Hub
}

func NewHandler(services *service.Services, repositories *repository.Repositories, db *mongo.Database, oauth *config.OauthConfig, auth *config.AuthConfig, i18n *config.I18nConfig, imageConfig *config.IImageConfig, hub service.Hub) *HandlerV1 {
	return &HandlerV1{
		repositories: repositories,
		db:           db,
		Services:     services,
		oauth:        *oauth,
		auth:         *auth,
		i18n:         *i18n,
		imageConfig:  *imageConfig,
		hub:          hub,
	}
}

func (h *HandlerV1) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{

		h.registerAuth(v1)
		oauth := v1.Group("/oauth")
		h.registerVkOAuth(oauth)
		h.registerGoogleOAuth(oauth)

		h.registerCategory(v1)
		h.RegisterLang(v1)
		h.RegisterCurrency(v1)
		h.RegisterCountry(v1)
		h.registerTag(v1)
		h.registerQuestion(v1)

		authenticated := v1.Group("", h.SetUserFromRequest)
		{
			h.registerAction(authenticated)
			h.registerAddress(authenticated)
			h.RegisterImage(authenticated)
			h.registerGql(authenticated)
			h.registerSubscribe(authenticated)
			h.registerProduct(authenticated)
			h.registerWs(authenticated)
			// h.registerNodeVote(authenticated)
			h.registerMessage(authenticated)
			h.registerMessageRoom(authenticated)
			h.registerOffer(authenticated)
			// h.registerNodedataVote(authenticated)
			h.registerReview(authenticated)
			h.registerTicket(authenticated)
			h.registerTrack(authenticated)
			h.RegisterUser(authenticated)
		}

		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "v1",
			})
		})
	}
}
