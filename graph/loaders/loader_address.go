package loaders

// import graph gophers with your other imports
import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *loaderReader) getAddress(ctx context.Context, addressIds []string) []*dataloader.Result[*model.Address] {
	gc, err := utils.GinContextFromContext(ctx)
	lang := gc.MustGet("i18nLocale").(string)
	if err != nil {
		panic(err)
	}

	keyOrder := make(map[string]int, len(addressIds))
	// var itemIDs []string
	for ix, key := range addressIds {
		// itemIDs = append(itemIDs, key)
		keyOrder[key] = ix
	}
	results := make([]*dataloader.Result[*model.Address], len(addressIds))

	dbRecords, err := u.Repo.Address.GqlGetAdresses(domain.RequestParams{
		Filter: bson.M{"osm_id": bson.M{"$in": addressIds}},
		Lang:   lang,
	})
	if err != nil {
		return results
	}

	// if DB error, return
	if err != nil {
		for i := 0; i < len(results); i++ {
			results[i] = &dataloader.Result[*model.Address]{Error: err}
		}
		return results
	}
	// enumerate records, put into output
	for _, record := range dbRecords {
		keyRecord := record.OsmID
		ix, ok := keyOrder[keyRecord]
		// if found, remove from index lookup map so we know elements were found
		if ok {
			results[ix] = &dataloader.Result[*model.Address]{Data: record, Error: nil}
			delete(keyOrder, keyRecord)
		}
	}
	// fill array positions with errors where not found in DB
	for addressID, ix := range keyOrder {
		err := fmt.Errorf("address not found %s", addressID)
		results[ix] = &dataloader.Result[*model.Address]{Data: nil, Error: err}
	}

	return results
}

func GetAddress(ctx context.Context, addressID string) (*model.Address, error) {
	loaders := For(ctx)
	return loaders.AddressLoader.Load(ctx, addressID)()
}

func GetAddresses(ctx context.Context, addressIDs []string) ([]*model.Address, []error) {
	loaders := For(ctx)
	return loaders.AddressLoader.LoadMany(ctx, addressIDs)()
}
