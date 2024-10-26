package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewReviewMongo(db *mongo.Database, i18n config.I18nConfig) *ReviewMongo {
	return &ReviewMongo{db: db, i18n: i18n}
}

func (r *ReviewMongo) FindReview(params domain.RequestParams) (domain.Response[model.Review], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Review
	var response domain.Response[model.Review]
	// var response domain.Response[model.Review]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Review]{}, err
	// }

	// cursor, err := r.db.Collection(TblReview).Find(ctx, filter, opts)
	// if err != nil {
	// 	return response, err
	// }
	// defer cursor.Close(ctx)

	// if er := cursor.All(ctx, &results); er != nil {
	// 	return response, er
	// }

	// resultSlice := make([]model.Review, len(results))
	// // for i, d := range results {
	// // 	resultSlice[i] = d
	// // }
	// copy(resultSlice, results)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Review]{}, err
	}
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		// "localField":   "user_id",
		// "foreignField": "_id",
		"let": bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": tblImage,
					"as":   "images",
					// "localField":   "_id",
					// "foreignField": "service_id",
					"let": bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})
	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	// cursor, err := r.db.Collection(TblNode).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	count, err := r.db.Collection(TblReview).CountDocuments(ctx, params.Filter)
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Review]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  results,
	}
	return response, nil
}

func (r *ReviewMongo) GetAllReview(params domain.RequestParams) (domain.Response[model.Review], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Review
	var response domain.Response[model.Review]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Review]{}, err
	}

	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Review, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblReview).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Review]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ReviewMongo) GqlGetReviews(params domain.RequestParams) ([]*model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Review
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usera",
		"let":  bson.D{{Key: "userId", Value: "$user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": tblImage,
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"user": bson.M{"$first": "$usera"}}}})

	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Review, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	// count, err := r.db.Collection(TblReview).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return results, err
	// }

	// results = []*model.Review{
	// 	Total: int(count),
	// 	Skip:  int(params.Options.Skip),
	// 	Limit: int(params.Options.Limit),
	// 	Data:  resultSlice,
	// }
	return results, nil
}

func (r *ReviewMongo) GqlGetCountReviews(params domain.RequestParams) (*model.ReviewInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results model.ReviewInfo
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return &results, err
	}

	var allItems []*model.Review
	cursor, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return &results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &allItems); er != nil {
		return &results, er
	}

	count := len(allItems)
	results.Count = &count

	var sum = int(0)
	for _, t := range allItems {
		sum += t.Rate
	}
	results.Value = &sum

	pipe = append(pipe,
		bson.D{
			{"$group", bson.D{
				{"_id", "$rate"},
				// {"average_price", bson.D{{"$avg", "$price"}}},
				{"count", bson.D{{"$sum", 1}}},
			}}})

	var allItemsGroup []*map[string]int
	cursorGroup, err := r.db.Collection(TblReview).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return &results, err
	}
	defer cursorGroup.Close(ctx)

	if er := cursorGroup.All(ctx, &allItemsGroup); er != nil {
		return &results, er
	}

	results.Ratings = allItemsGroup

	// if typeC, ok := f["type"]; ok {
	// 	typeContent = typeC.(string)
	// }
	// keys := make([]int, 0, len(allItemsGroup))
	// counts := make([]int, 0, len(allItemsGroup))
	// for k := range allItemsGroup {
	// 	f := *allItemsGroup[k]
	// 	counts = append(counts, f["count"])
	// 	keys = append(keys, f["_id"])
	// }
	// sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	// for _, v := range keys {
	// 	f := *allItemsGroup[counts[v-1]]
	// 	fmt.Println("allItemsGroup:::", v, f["_id"], "-", f["count"])
	// }

	return &results, nil
}

func (r *ReviewMongo) CreateReview(userID string, review *model.Review) (*model.Review, error) {
	var result *model.Review

	collection := r.db.Collection(TblReview)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// var existReview model.Review
	// r.db.Collection(TblReview).FindOne(ctx, bson.M{"node_id": review.NodeID, "user_id": userIDPrimitive}).Decode(&existReview)

	// if existReview.NodeID.IsZero() {
	updatedAt := review.UpdatedAt
	if updatedAt.IsZero() {
		updatedAt = time.Now()
	}

	newReview := model.ReviewInput{
		Review:    review.Review,
		Rate:      review.Rate,
		NodeID:    review.NodeID,
		UserID:    userIDPrimitive,
		CreatedAt: updatedAt,
		UpdatedAt: updatedAt,
	}

	res, err := collection.InsertOne(ctx, newReview)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblReview).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	// } else {
	// 	updatedAt := review.UpdatedAt
	// 	if updatedAt.IsZero() {
	// 		updatedAt = time.Now()
	// 	}

	// 	updateReview := &model.ReviewInput{
	// 		Rate:      review.Rate,
	// 		Review:    review.Review,
	// 		UpdatedAt: updatedAt,
	// 	}
	// 	result, err = r.UpdateReview(existReview.ID.Hex(), userID, updateReview)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	return result, nil
}

func (r *ReviewMongo) UpdateReview(id string, userID string, data *model.ReviewInput) (*model.Review, error) {
	var result *model.Review
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblReview)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Rate != 0 {
		newData["rate"] = data.Rate
	}
	if data.Review != "" {
		newData["review"] = data.Review
	}
	newData["updated_at"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ReviewMongo) DeleteReview(id string) (*model.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = &model.Review{}
	collection := r.db.Collection(TblReview)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}
