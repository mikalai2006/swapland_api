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

type AddressMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewAddressMongo(db *mongo.Database, i18n config.I18nConfig) *AddressMongo {
	return &AddressMongo{db: db, i18n: i18n}
}

func (r *AddressMongo) FindAddress(params domain.RequestParams) (domain.Response[domain.Address], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Address
	var response domain.Response[domain.Address]
	filter, opts, err := CreateFilterAndOptions(params)
	if err != nil {
		return domain.Response[domain.Address]{}, err
	}

	cursor, err := r.db.Collection(TblAddress).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Address, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAddress).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Address]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *AddressMongo) GetAllAddress(params domain.RequestParams) (domain.Response[domain.Address], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Address
	var response domain.Response[domain.Address]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Address]{}, err
	}

	cursor, err := r.db.Collection(TblAddress).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Address, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAddress).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Address]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *AddressMongo) CreateAddress(userID string, address *domain.AddressInput) (*domain.Address, error) {
	var result *domain.Address

	collection := r.db.Collection(TblAddress)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newAddress := domain.Address{
		UserID:    userIDPrimitive,
		Address:   address.Address,
		OsmID:     address.OsmID,
		Lang:      address.Lang,
		DAddress:  address.DAddress,
		Props:     address.Props,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newAddress)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblAddress).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AddressMongo) DeleteAddress(id string) (model.Address, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Address{}
	collection := r.db.Collection(TblAddress)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	// if err != nil {
	// 	return result, err
	// }

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *AddressMongo) GqlGetAdresses(params domain.RequestParams) ([]*model.Address, error) {
	// fmt.Println("GqlGetAdresses")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Address
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	cursor, err := r.db.Collection(TblAddress).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Address, len(results))
	copy(resultSlice, results)

	return results, nil
}
