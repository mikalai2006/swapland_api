package loaders

// import graph gophers with your other imports
import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader/v7"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// tagReader reads Tags from a database
type loaderReader struct {
	DB   *mongo.Database
	Repo *repository.Repositories
}

// getTags implements a batch function that can retrieve many tags by ID,
// for use in a dataloader
func (u *loaderReader) getTags(ctx context.Context, tagIds []string) []*dataloader.Result[*model.Tag] {
	gc, err := utils.GinContextFromContext(ctx)
	lang := gc.MustGet("i18nLocale").(string)
	if err != nil {
		panic(err)
	}
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(tagIds))
	// collect the keys to search for
	// var tagIDs []string
	for ix, key := range tagIds {
		// tagIDs = append(tagIDs, key)
		keyOrder[key] = ix
	}
	// construct an output array of dataloader results
	results := make([]*dataloader.Result[*model.Tag], len(tagIds))

	listIDs := []primitive.ObjectID{}
	for i := range tagIds {
		uIDPrimitive, err := primitive.ObjectIDFromHex(tagIds[i])
		if err != nil {
			panic(err)
		}
		listIDs = append(listIDs, uIDPrimitive)
	}

	dbRecords, err := u.Repo.Tag.GqlGetTags(domain.RequestParams{
		// Options: domain.Options{Limit: int64(*limit)},
		Filter: bson.M{"_id": bson.M{"$in": listIDs}},
		Lang:   lang,
	})
	if err != nil {
		return results
	}

	// if DB error, return
	if err != nil {
		for i := 0; i < len(results); i++ {
			results[i] = &dataloader.Result[*model.Tag]{Error: err}
		}
		return results
	}
	// enumerate records, put into output
	for _, record := range dbRecords {
		keyRecord := record.ID.Hex()
		ix, ok := keyOrder[keyRecord]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result[*model.Tag]{Data: record, Error: nil}
			delete(keyOrder, keyRecord)
		}
	}
	// fill array positions with errors where not found in DB
	for tagID, ix := range keyOrder {
		err := fmt.Errorf("tag not found %s", tagID)
		results[ix] = &dataloader.Result[*model.Tag]{Data: nil, Error: err}
	}

	// fmt.Println("result tags=", len(result), tagIds)

	return results
}

// handleError creates array of result with the same error repeated for as many items requested
func handleError[T any](itemsLength int, err error) []*dataloader.Result[T] {
	result := make([]*dataloader.Result[T], itemsLength)
	for i := 0; i < itemsLength; i++ {
		result[i] = &dataloader.Result[T]{Error: err}
	}
	return result
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	AddressLoader *dataloader.Loader[string, *model.Address]
	TagsLoader    *dataloader.Loader[string, *model.Tag]
	TagoptLoader  *dataloader.Loader[string, []*model.Question]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(mongoDB *mongo.Database, repositories *repository.Repositories) *Loaders {
	// define the data loader
	ur := &loaderReader{DB: mongoDB, Repo: repositories}
	// dataloader.WithClearCacheOnBatch[string, *model.Tag]()
	// dataloader.WithWait[string, *model.Address](time.Millisecond)
	//dataloader.WithCache[string, *model.Tagopt](&dataloader.NoCache[string, *model.Tagopt]{})
	return &Loaders{
		AddressLoader: dataloader.NewBatchedLoader(ur.getAddress, dataloader.WithClearCacheOnBatch[string, *model.Address]()),
		TagsLoader:    dataloader.NewBatchedLoader(ur.getTags, dataloader.WithWait[string, *model.Tag](time.Millisecond)),
		// TagoptLoader:  dataloader.NewBatchedLoader(ur.getTagopt, dataloader.WithClearCacheOnBatch[string, []*model.Question]()),
	}
}

// Middleware injects data loaders into the context
func Middleware(mongoDB *mongo.Database, repositories *repository.Repositories) gin.HandlerFunc {
	// return a middleware that injects the loader to the request context
	return func(c *gin.Context) {
		loader := NewLoaders(mongoDB, repositories)
		ctx := context.WithValue(c.Request.Context(), loadersKey, loader)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	// fmt.Println("For")
	return ctx.Value(loadersKey).(*Loaders)
}

// GetTag returns single tag by id efficiently
func GetTag(ctx context.Context, tagID string) (*model.Tag, error) {
	loaders := For(ctx)
	return loaders.TagsLoader.Load(ctx, tagID)()
}

// GetTags returns many tags by ids efficiently
func GetTags(ctx context.Context, tagIDs []string) ([]*model.Tag, []error) {
	// fmt.Println("GetTags1", tagIDs)
	loaders := For(ctx)
	// fmt.Println("GetTags2")
	return loaders.TagsLoader.LoadMany(ctx, tagIDs)()
}
