package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ImageMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewImageMongo(db *mongo.Database, i18n config.I18nConfig) *ImageMongo {
	return &ImageMongo{db: db, i18n: i18n}
}

func (r *ImageMongo) CreateImage(userID string, data *model.ImageInput) (model.Image, error) {
	var result model.Image

	collection := r.db.Collection(tblImage)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	// var ServiceID primitive.ObjectID
	// if data.ServiceID != "" {
	// 	ServiceID, err = primitive.ObjectIDFromHex(data.ServiceID)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// } else {
	// 	ServiceID = primitive.NilObjectID
	// }

	newImage := model.ImageInputMongo{
		UserID:      userIDPrimitive,
		Service:     data.Service,
		ServiceID:   data.ServiceID,
		Path:        data.Path,
		Title:       data.Title,
		Ext:         data.Ext,
		Dir:         data.Dir,
		Description: data.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// if primitive.IsValidObjectID(data.ServiceID) {
	// 	serviceIDPrimitive, err := primitive.ObjectIDFromHex(data.ServiceID)
	// 	if err != nil {
	// 		return result, err
	// 	}
	// 	newImage.ServiceID = serviceIDPrimitive
	// 	fmt.Println("valid serviceId")
	// }

	res, err := collection.InsertOne(ctx, newImage)
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ImageMongo) GetImage(id string) (model.Image, error) {
	var result model.Image

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Image{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblImage).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return model.Image{}, err
	}

	return result, nil
}

func (r *ImageMongo) GetImageDirs(id string) ([]interface{}, error) {
	var result []interface{}

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"user_id": userIDPrimitive}
	// pipe := mongo.Pipeline{}

	// pipe = append(pipe, bson.D{{"$match", bson.M{"user_id": userIDPrimitive}}})
	// pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
	// 	"from":         "component_presets",
	// 	"as":           "presets",
	// 	"localField":   "_id",
	// 	"foreignField": "component_id",
	// }}})

	result, err = r.db.Collection(tblImage).Distinct(ctx, "dir", filter) //.Aggregate(ctx, pipe) // (ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *ImageMongo) FindImage(params domain.RequestParams) (domain.Response[model.Image], error) {
	var results []model.Image
	var response domain.Response[model.Image]

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblImage)

	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Image]{}, err
	}

	cursor, err := collection.Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Image, len(results))
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	// options.SetSkip(params.Skip)
	count, err := collection.CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Image]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *ImageMongo) DeleteImage(id string) (model.Image, error) {
	var result = model.Image{}

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblImage)

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

func (r *ImageMongo) GqlGetImages(params domain.RequestParams) ([]*model.Image, error) {
	fmt.Println("GqlGetImages")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Image
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(tblImage).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Image, len(results))

	copy(resultSlice, results)
	return results, nil
}
