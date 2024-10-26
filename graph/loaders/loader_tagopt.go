package loaders

// // import graph gophers with your other imports
// import (
// 	"context"
// 	"fmt"
// 	"strings"

// 	"github.com/graph-gophers/dataloader/v7"
// 	"github.com/mikalai2006/swapland-api/graph/model"
// 	"github.com/mikalai2006/swapland-api/internal/domain"
// 	"github.com/mikalai2006/swapland-api/internal/utils"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func (u *loaderReader) getTagopt(ctx context.Context, tagoptIds []string) []*dataloader.Result[[]*model.Question] {
// 	gc, err := utils.GinContextFromContext(ctx)
// 	lang := gc.MustGet("i18nLocale").(string)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Println("tagoptIds: ", tagoptIds, len(tagoptIds))

// 	keyOrder := make(map[string]int, len(tagoptIds))

// 	query := []bson.M{}
// 	for ix, key := range tagoptIds {
// 		// itemIDs = append(itemIDs, key)

// 		idITems := strings.Split(key, ".")

// 		itemID, err := primitive.ObjectIDFromHex(idITems[1])
// 		if err != nil {
// 			panic(err)
// 		}
// 		query = append(query, bson.M{"$and": []bson.M{{"tag_id": itemID}, {"osm_id": idITems[0]}}})
// 		keyOrder[key] = ix
// 	}
// 	results := make([]*dataloader.Result[[]*model.Question], len(tagoptIds))
// 	// fmt.Println("query=", query)
// 	dbRecords, err := u.Repo.Tagopt.GqlGetTagopts(domain.RequestParams{
// 		Filter: bson.D{{"$or", query}},
// 		Lang:   lang,
// 	})
// 	if err != nil {
// 		return results
// 	}
// 	fmt.Println("dbRecords=", len(dbRecords))

// 	// if DB error, return
// 	if err != nil {
// 		for i := 0; i < len(results); i++ {
// 			results[i] = &dataloader.Result[[]*model.Question]{Error: err}
// 		}
// 		return results
// 	}

// 	// // prepairs result for output
// 	// res := map[string][]*model.Tagopt{}
// 	// for _, record := range dbRecords {
// 	// 	keyRecord := fmt.Sprintf("%v.%v", record.OsmID, record.TagID.Hex())
// 	// 	res[keyRecord] = append(res[keyRecord], record)
// 	// 	_, ok := keyOrder[keyRecord]
// 	// 	// if found, remove from index lookup map so we know elements were found
// 	// 	if ok {
// 	// 		// results[ix] = &dataloader.Result[*model.Tagopt]{Data: record, Error: nil}
// 	// 		// delete(keyOrder, keyRecord)
// 	// 	}
// 	// }

// 	// // enumerate records, put into output
// 	// for _, record := range dbRecords {
// 	// 	keyRecord := fmt.Sprintf("%v.%v", record.OsmID, record.TagID.Hex())
// 	// 	ix, ok := keyOrder[keyRecord]

// 	// 	// get res
// 	// 	rec := []*model.Tagopt{}
// 	// 	_, okk := res[keyRecord]
// 	// 	if okk {
// 	// 		rec = res[keyRecord]
// 	// 	}

// 	// 	// if found, remove from index lookup map so we know elements were found
// 	// 	if ok {
// 	// 		results[ix] = &dataloader.Result[[]*model.Tagopt]{Data: rec, Error: nil}
// 	// 		delete(keyOrder, keyRecord)
// 	// 	}
// 	// }
// 	// fill array positions with errors where not found in DB
// 	for tagoptID, ix := range keyOrder {
// 		err := fmt.Errorf("tagopt not found %s", tagoptID)
// 		results[ix] = &dataloader.Result[[]*model.Question]{Data: []*model.Question{}, Error: err}
// 	}

// 	return results
// }

// func GetTagopt(ctx context.Context, itemID string) ([]*model.Question, error) {
// 	loaders := For(ctx)
// 	return loaders.TagoptLoader.Load(ctx, itemID)()
// }

// // func GetTagopts(ctx context.Context, itemIDs []string) ([]*model.Tagopt, []error) {
// // 	loaders := For(ctx)
// // 	return loaders.TagoptLoader.LoadMany(ctx, itemIDs)()
// // }
