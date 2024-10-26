package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubscribeMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewSubscribeMongo(db *mongo.Database, i18n config.I18nConfig) *SubscribeMongo {
	return &SubscribeMongo{db: db, i18n: i18n}
}

func (r *SubscribeMongo) FindSubscribe(params domain.RequestParams) (domain.Response[model.Subscribe], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Subscribe
	var response domain.Response[model.Subscribe]
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblSubscribe).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Subscribe, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(TblSubscribe).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Subscribe]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *SubscribeMongo) CreateSubscribe(userID string, Subscribe *model.SubscribeInput) (*model.Subscribe, error) {
	var result *model.Subscribe

	collection := r.db.Collection(TblSubscribe)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	// subUserIDPrimitive, err := primitive.ObjectIDFromHex(Subscribe.SubUserID)
	// if err != nil {
	// 	return nil, err
	// }

	newSubscribe := model.SubscribeInput{
		UserID:    userIDPrimitive,
		SubUserID: Subscribe.SubUserID,
		Status:    Subscribe.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newSubscribe)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblSubscribe).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SubscribeMongo) GqlGetSubscribes(params domain.RequestParams) ([]*model.Subscribe, error) {
	fmt.Println("GqlGetSubscribes")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Subscribe
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblSubscribe).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Subscribe, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *SubscribeMongo) GqlGetIamSubscribe(userID string, nodeID string) (*model.Subscribe, error) {
	fmt.Println("GqlGetIamSubscribe")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result *model.Subscribe

	nodeIDPrimitive, err := primitive.ObjectIDFromHex(nodeID)
	if err != nil {
		return result, err
	}
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	if err := r.db.Collection(TblSubscribe).FindOne(ctx, bson.D{{"node_id", nodeIDPrimitive}, {"user_id", userIDPrimitive}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return result, model.ErrSubscribeNotFound
		}
		return result, err
	}
	return result, nil
}

func (r *SubscribeMongo) UpdateSubscribe(id string, userID string, data *model.Subscribe) (*model.Subscribe, error) {
	var result *model.Subscribe
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblSubscribe)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
		"status":     data.Status,
		"updated_at": time.Now(),
	}})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *SubscribeMongo) DeleteSubscribe(id string) (model.Subscribe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Subscribe{}
	collection := r.db.Collection(TblSubscribe)

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
