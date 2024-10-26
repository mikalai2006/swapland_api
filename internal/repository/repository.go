package repository

import (
	"reflect"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Action interface {
	FindAction(params domain.RequestParams) (domain.Response[model.Action], error)
	GetAllAction(params domain.RequestParams) (domain.Response[model.Action], error)
	CreateAction(userID string, tag *model.ActionInput) (*model.Action, error)
	UpdateAction(id string, userID string, data *model.ActionInput) (*model.Action, error)
	DeleteAction(id string) (model.Action, error)
	GqlGetActions(params domain.RequestParams) ([]*model.Action, error)
}

type Address interface {
	FindAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	GetAllAddress(params domain.RequestParams) (domain.Response[domain.Address], error)
	CreateAddress(userID string, address *domain.AddressInput) (*domain.Address, error)
	DeleteAddress(id string) (model.Address, error)
	GqlGetAdresses(params domain.RequestParams) ([]*model.Address, error)
}

type Authorization interface {
	CreateAuth(auth *domain.SignInInput) (string, error)
	GetAuth(id string) (domain.Auth, error)
	CheckExistAuth(auth *domain.SignInInput) (domain.Auth, error)
	GetByCredentials(auth *domain.SignInInput) (domain.Auth, error)
	SetSession(authID primitive.ObjectID, session domain.Session) error
	VerificationCode(userID string, code string) error
	RefreshToken(refreshToken string) (domain.Auth, error)
	RemoveRefreshToken(refreshToken string) (string, error)
}

type Track interface {
	FindTrack(params domain.RequestParams) (domain.Response[domain.Track], error)
	GetAllTrack(params domain.RequestParams) (domain.Response[domain.Track], error)
	CreateTrack(userID string, track *domain.Track) (*domain.Track, error)
}

type Product interface {
	FindProduct(params *model.ProductFilter) (domain.Response[model.Product], error)
	CreateProduct(userID string, product *model.ProductInputData) (*model.Product, error)
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
	FindOffer(params *model.OfferFilter) (domain.Response[model.Offer], error)
	GetOffer(id string) (*model.Offer, error)
	CreateOffer(userID string, data *model.OfferInput) (*model.Offer, error)
	UpdateOffer(id string, userID string, data *model.Offer) (*model.Offer, error)
	DeleteOffer(id string) (model.Offer, error)
	GqlGetOffers(params domain.RequestParams) ([]*model.Offer, error)
}

type Review interface {
	FindReview(params domain.RequestParams) (domain.Response[model.Review], error)
	GetAllReview(params domain.RequestParams) (domain.Response[model.Review], error)
	CreateReview(userID string, review *model.Review) (*model.Review, error)
	UpdateReview(id string, userID string, data *model.ReviewInput) (*model.Review, error)
	DeleteReview(id string) (*model.Review, error)

	GqlGetReviews(params domain.RequestParams) ([]*model.Review, error)
	GqlGetCountReviews(params domain.RequestParams) (*model.ReviewInfo, error)
}

type User interface {
	GetUser(id string) (model.User, error)
	FindUser(params domain.RequestParams) (domain.Response[model.User], error)
	CreateUser(userID string, user *model.User) (*model.User, error)
	DeleteUser(id string) (model.User, error)
	UpdateUser(id string, user *model.User) (model.User, error)
	Iam(userID string) (model.User, error)

	SetStat(userID string, statData model.UserStat) (model.User, error)
	SetBal(userID string, value int) (model.User, error)
	GqlGetUsers(params domain.RequestParams) ([]*model.User, error)
}

type Image interface {
	CreateImage(userID string, data *model.ImageInput) (model.Image, error)
	GetImage(id string) (model.Image, error)
	GetImageDirs(id string) ([]interface{}, error)
	FindImage(params domain.RequestParams) (domain.Response[model.Image], error)
	DeleteImage(id string) (model.Image, error)

	GqlGetImages(params domain.RequestParams) ([]*model.Image, error)
}

type Lang interface {
	CreateLanguage(userID string, data *domain.LanguageInput) (domain.Language, error)
	GetLanguage(id string) (domain.Language, error)
	FindLanguage(params domain.RequestParams) (domain.Response[domain.Language], error)
	UpdateLanguage(id string, data interface{}) (domain.Language, error)
	DeleteLanguage(id string) (domain.Language, error)
}

type Currency interface {
	CreateCurrency(userID string, data *domain.CurrencyInput) (domain.Currency, error)
	GetCurrency(id string) (domain.Currency, error)
	FindCurrency(params domain.RequestParams) (domain.Response[domain.Currency], error)
	UpdateCurrency(id string, data interface{}) (domain.Currency, error)
	DeleteCurrency(id string) (domain.Currency, error)
}

type Country interface {
	CreateCountry(userID string, data *domain.CountryInput) (domain.Country, error)
	GetCountry(id string) (domain.Country, error)
	FindCountry(params domain.RequestParams) (domain.Response[domain.Country], error)
	UpdateCountry(id string, data interface{}) (domain.Country, error)
	DeleteCountry(id string) (domain.Country, error)
}

type Tag interface {
	FindTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	GetAllTag(params domain.RequestParams) (domain.Response[model.Tag], error)
	CreateTag(userID string, tag *model.Tag) (*model.Tag, error)
	UpdateTag(id string, userID string, data *model.Tag) (*model.Tag, error)
	DeleteTag(id string) (model.Tag, error)
	GqlGetTags(params domain.RequestParams) ([]*model.Tag, error)
}

type Question interface {
	FindQuestion(params *model.QuestionFilter) (domain.Response[model.Question], error)
	CreateQuestion(userID string, data *model.QuestionInput) (*model.Question, error)
	UpdateQuestion(id string, userID string, data *model.QuestionInput) (*model.Question, error)
	DeleteQuestion(id string) (model.Question, error)
}

type Ticket interface {
	FindTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	GetAllTicket(params domain.RequestParams) (domain.Response[model.Ticket], error)
	CreateTicket(userID string, ticket *model.Ticket) (*model.Ticket, error)
	CreateTicketMessage(userID string, message *model.TicketMessage) (*model.TicketMessage, error)
	DeleteTicket(id string) (model.Ticket, error)
	GqlGetTickets(params domain.RequestParams) ([]*model.Ticket, error)
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
	GqlGetCategorys(params domain.RequestParams) ([]*model.Category, error)
}

type Repositories struct {
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

func NewRepositories(mongodb *mongo.Database, i18n config.I18nConfig) *Repositories {
	return &Repositories{
		Action:        NewActionMongo(mongodb, i18n),
		Address:       NewAddressMongo(mongodb, i18n),
		Category:      NewCategoryMongo(mongodb, i18n),
		Authorization: NewAuthMongo(mongodb),
		Lang:          NewLangMongo(mongodb, i18n),
		Country:       NewCountryMongo(mongodb, i18n),
		Currency:      NewCurrencyMongo(mongodb, i18n),
		Image:         NewImageMongo(mongodb, i18n),
		Review:        NewReviewMongo(mongodb, i18n),
		User:          NewUserMongo(mongodb, i18n),
		Track:         NewTrackMongo(mongodb, i18n),
		Product:       NewProductMongo(mongodb, i18n),
		NodeVote:      NewNodeVoteMongo(mongodb, i18n),
		Message:       NewMessageMongo(mongodb, i18n),
		MessageRoom:   NewMessageRoomMongo(mongodb, i18n),
		// MessageImage:  NewMessageImageMongo(mongodb, i18n),
		Offer:    NewOfferMongo(mongodb, i18n),
		Tag:      NewTagMongo(mongodb, i18n),
		Question: NewQuestionMongo(mongodb, i18n),
		Ticket:   NewTicketMongo(mongodb, i18n),

		Subscribe: NewSubscribeMongo(mongodb, i18n),
	}
}

// func getPaginationOpts(pagination *domain.PaginationQuery) *options.FindOptions {
// 	var opts *options.FindOptions
// 	if pagination != nil {
// 		opts = &options.FindOptions{
// 			Skip:  pagination.GetSkip(),
// 			Limit: pagination.GetLimit(),
// 		}
// 	}

// 	return opts
// }

func createFilter[V any](filterData *V) any {
	var filter V

	filterReflect := reflect.ValueOf(filterData)
	// fmt.Println("========== filterReflect ===========")
	// fmt.Println("struct > ", filterReflect)
	// fmt.Println("struct type > ", filterReflect.Type())
	filterIndirectData := reflect.Indirect(filterReflect)
	// fmt.Println("filter data > ", filterIndirectData)
	// fmt.Println("filter numField > ", filterIndirectData.NumField())
	dataFilter := bson.M{}

	var tagJSON, tagPrimitive string
	for i := 0; i < filterIndirectData.NumField(); i++ {
		field := filterIndirectData.Field(i)
		if field.Kind() == reflect.Ptr {
			field = reflect.Indirect(field)
		}
		typeField := filterIndirectData.Type().Field(i)
		tag := typeField.Tag
		// tagBson = tag.Get("bson")
		tagJSON = tag.Get("json")
		tagPrimitive = tag.Get("primitive")
		switch field.Kind() {
		case reflect.String:
			value := field.String()
			if tagPrimitive == "true" {
				id, _ := primitive.ObjectIDFromHex(value)
				// fmt.Println("===== string add ", tag, value)
				dataFilter[tagJSON] = id
			} else {
				dataFilter[tagJSON] = value
			}

		case reflect.Bool:
			value := field.Bool()
			dataFilter[tagJSON] = value

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value := field.Int()
			dataFilter[tagJSON] = value

		default:

		}

		// fmt.Println(tagBson, tagJSON, tagPrimitive, fmt.Sprintf("[%s]", field), field.Kind(), field)
	}

	// structure := reflect.ValueOf(&filter)
	// fmt.Println("========== filter ===========")
	// fmt.Println("struct > ", structure)
	// fmt.Println("struct type > ", structure.Type())
	// fmt.Println("filter data > ", reflect.Indirect(structure))
	// fmt.Println("filter numField > ", reflect.Indirect(structure).NumField())

	// fmt.Println("========== result ===========")
	// fmt.Println("dataFilter > ", dataFilter)
	return filter
}
