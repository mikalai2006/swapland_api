package service

import (
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"github.com/mikalai2006/swapland-api/pkg/auths"
	"github.com/mikalai2006/swapland-api/pkg/hasher"
)

type Action interface {
	FindAction(params domain.RequestParams) (domain.Response[model.Action], error)
	GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error)
	CreateAction(userID string, data *model.ActionInput) (*model.Action, error)
	UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error)
	DeleteAction(id string) (model.Action, error)
}
type Address interface {
	CreateAddress(userID string, address domain.ResponseNominatim) (*model.Address, error)
	FindAddress(input *model.AddressFilter) (domain.Response[model.Address], error)
	// GetAllAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	UpdateAddress(id string, userID string, data domain.ResponseNominatim) (*model.Address, error)
	DeleteAddress(id string) (model.Address, error)
}

type Authorization interface {
	CreateAuth(auth *domain.SignInInput) (string, error)
	GetAuth(id string) (domain.Auth, error)
	SignIn(input *domain.SignInInput) (domain.ResponseTokens, error)
	ExistAuth(auth *domain.SignInInput) (domain.Auth, error)
	CreateSession(auth *domain.Auth) (domain.ResponseTokens, error)
	VerificationCode(userID string, code string) error
	RefreshTokens(refreshToken string) (domain.ResponseTokens, error)
	RemoveRefreshTokens(refreshToken string) (string, error)
}

type Track interface {
	FindTrack(params domain.RequestParams) (domain.Response[domain.Track], error)
	GetAllTrack(params domain.RequestParams) (domain.Response[domain.Track], error)
	CreateTrack(userID string, track *domain.Track) (*domain.Track, error)
}

type Product interface {
	FindProduct(params *model.ProductFilter) (domain.Response[model.Product], error)
	CreateProduct(userID string, node *model.ProductInputData) (*model.Product, error)
	UpdateProduct(id string, userID string, data *model.Product) (*model.Product, error)
	DeleteProduct(id string) (model.Product, error)
}

type Message interface {
	CreateMessage(userID string, message *model.MessageInput) (*model.Message, error)
	FindMessage(params *model.MessageFilter) (domain.Response[model.Message], error)
	UpdateMessage(id string, userID string, data *model.MessageInput) (*model.Message, error)
	DeleteMessage(id string) (model.Message, error)
	GetGroupForUser(userID string) ([]model.MessageGroupForUser, error)
}

type MessageRoom interface {
	CreateMessageRoom(userID string, message *model.MessageRoom) (*model.MessageRoom, error)
	FindMessageRoom(params *model.MessageRoomFilter) (domain.Response[model.MessageRoom], error)
	UpdateMessageRoom(id string, userID string, data *model.MessageRoom) (*model.MessageRoom, error)
	DeleteMessageRoom(id string) (model.MessageRoom, error)
	// GetGroupForUser(userID string) ([]model.MessageGroupForUser, error)
}

// type MessageImage interface {
// 	CreateMessageImage(userID string, data *model.MessageImageInput) (model.MessageImage, error)
// 	GetMessageImage(id string) (model.MessageImage, error)
// 	DeleteMessageImage(id string) (model.MessageImage, error)
// }

type NodeVote interface {
	CreateNodeVote(userID string, data *model.NodeVote) (*model.NodeVote, error)
	FindNodeVote(params domain.RequestParams) (domain.Response[model.NodeVote], error)
	UpdateNodeVote(id string, userID string, data *model.NodeVoteInput) (*model.NodeVote, error)
	DeleteNodeVote(id string) (model.NodeVote, error)
}

type Offer interface {
	CreateOffer(userID string, data *model.OfferInput) (*model.Offer, error)
	FindOffer(params *model.OfferFilter) (domain.Response[model.Offer], error)
	UpdateOffer(id string, userID string, data *model.Offer) (*model.Offer, error)
	DeleteOffer(id string) (model.Offer, error)
}

type Review interface {
	CreateReview(userID string, review *model.Review) (*model.Review, error)
	FindReview(params domain.RequestParams) (domain.Response[model.Review], error)
	UpdateReview(id string, userID string, data *model.ReviewInput) (*model.Review, error)
	DeleteReview(id string) (*model.Review, error)

	GetAllReview(params domain.RequestParams) (domain.Response[model.Review], error)
}
type User interface {
	GetUser(id string) (model.User, error)
	FindUser(params domain.RequestParams) (domain.Response[model.User], error)
	CreateUser(userID string, user *model.User) (*model.User, error)
	DeleteUser(id string) (model.User, error)
	UpdateUser(id string, user *model.User) (model.User, error)
	Iam(userID string) (model.User, error)

	SetStat(id string, statData model.UserStat) (model.User, error)
	SetBal(userID string, value int) (model.User, error)
}

type Image interface {
	CreateImage(userID string, data *model.ImageInput) (model.Image, error)
	GetImage(id string) (model.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[model.Image], error)
	DeleteImage(id string) (model.Image, error)
}
type Country interface {
	CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error)
	GetCountry(id string) (domain.Country, error)
	FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error)
	UpdateCountry(id string, data interface{}) (domain.Country, error)
	DeleteCountry(id string) (domain.Country, error)
}

type Currency interface {
	CreateCurrency(userID string, data *domain.CurrencyInput) (domain.Currency, error)
	GetCurrency(id string) (domain.Currency, error)
	FindCurrency(params domain.RequestParams) (domain.Response[domain.Currency], error)
	UpdateCurrency(id string, data interface{}) (domain.Currency, error)
	DeleteCurrency(id string) (domain.Currency, error)
}

type Lang interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Tag interface {
	FindTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	GetAllTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	CreateTag(userID string, tag *model.Tag) (*model.Tag, error)
	UpdateTag(id string, userID string, data *model.Tag) (*model.Tag, error)
	DeleteTag(id string) (model.Tag, error)
}
type Question interface {
	FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error)
	CreateQuestion(userID string, question *model.QuestionInput) (*model.Question, error)
	UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error)
	DeleteQuestion(id string) (model.Question, error)
}
type Ticket interface {
	FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error)
	CreateTicketMessage(userID string, message *model.TicketMessage) (*model.TicketMessage, error)
	DeleteTicket(id string) (model.Ticket, error)
}

type Subscribe interface {
	FindSubscribe(params domain.RequestParams) (domain.Response[model.Subscribe], error)
	CreateSubscribe(userID string, data *model.SubscribeInput) (*model.Subscribe, error)
	UpdateSubscribe(id string, userID string, data *model.Subscribe) (*model.Subscribe, error)
	DeleteSubscribe(id string) (model.Subscribe, error)
}
type Category interface {
	FindCategory(params domain.RequestParams) (domain.Response[model.Category], error)
	GetAllCategory(params domain.RequestParams) (domain.Response[model.Category], error)
	CreateCategory(userID string, Category *model.Category) (*model.Category, error)
	UpdateCategory(id string, userID string, data *model.CategoryInput) (*model.Category, error)
	DeleteCategory(id string) (model.Category, error)
}

type Services struct {
	Action
	Address
	Category
	Authorization
	Lang
	Country
	Currency
	Image
	Review
	User
	Track
	Product
	Message
	MessageRoom
	// MessageImage
	NodeVote
	Offer
	Tag
	Question
	Ticket

	Subscribe
}

type ConfigServices struct {
	Repositories           *repository.Repositories
	Hasher                 hasher.PasswordHasher
	TokenManager           auths.TokenManager
	OtpGenerator           utils.Generator
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeLength int
	I18n                   config.I18nConfig
	ImageService           config.IImageConfig
	Hub                    *Hub
}

func NewServices(cfgService *ConfigServices) *Services {
	User := NewUserService(cfgService.Repositories.User, cfgService.Hub)
	Authorization := NewAuthService(
		cfgService.Repositories.Authorization,
		cfgService.Hasher,
		cfgService.TokenManager,
		cfgService.RefreshTokenTTL,
		cfgService.AccessTokenTTL,
		cfgService.OtpGenerator,
		cfgService.VerificationCodeLength,
		User,
		cfgService.Hub,
	)
	Action := NewActionService(cfgService.Repositories.Action, cfgService.I18n)
	Address := NewAddressService(cfgService.Repositories.Address, cfgService.I18n)
	Category := NewCategoryService(cfgService.Repositories.Category, cfgService.I18n)
	// Review := NewReviewService(cfgService.Repositories.Review)
	Lang := NewLangService(cfgService.Repositories, cfgService.I18n)
	Country := NewCountryService(cfgService.Repositories, cfgService.I18n)
	Currency := NewCurrencyService(cfgService.Repositories, cfgService.I18n)
	Image := NewImageService(cfgService.Repositories.Image, cfgService.ImageService)
	Track := NewTrackService(cfgService.Repositories.Track)
	Product := NewProductService(cfgService.Repositories.Product, User, cfgService.Hub)
	MessageRoom := NewMessageRoomService(cfgService.Repositories.MessageRoom, cfgService.Hub)
	Message := NewMessageService(cfgService.Repositories.Message, cfgService.Hub, MessageRoom)
	// MessageImage := NewMessageImageService(cfgService.Repositories.MessageImage, cfgService.ImageService)
	NodeVote := NewNodeVoteService(cfgService.Repositories.NodeVote)
	// Nodedata:=     NewNodedataService(cfgService.Repositories.Nodedata, cfgService.Repositories.User, cfgService.Repositories.NodedataVote, *Services)
	Tag := NewTagService(cfgService.Repositories.Tag)
	Question := NewQuestionService(cfgService.Repositories.Question, cfgService.Hub)
	Ticket := NewTicketService(cfgService.Repositories.Ticket)
	Subscribe := NewSubscribeService(cfgService.Repositories.Subscribe)

	return &Services{
		Authorization: Authorization,
		Action:        Action,
		Address:       Address,
		Category:      Category,
		User:          User,
		Lang:          Lang,
		Subscribe:     Subscribe,
		Country:       Country,
		Image:         Image,
		Track:         Track,
		Product:       Product,
		Message:       Message,
		MessageRoom:   MessageRoom,
		// MessageImage:  MessageImage,
		NodeVote: NodeVote,
		Offer:    NewOfferService(cfgService.Repositories.Offer, User, cfgService.Hub, Message, MessageRoom),
		Tag:      Tag,
		Question: Question,
		Ticket:   Ticket,
		Review:   NewReviewService(cfgService.Repositories.Review, User),
		Currency: Currency,
	}
}
