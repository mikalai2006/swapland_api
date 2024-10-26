package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CurrencyMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewCurrencyMongo(db *mongo.Database, i18n config.I18nConfig) *CurrencyMongo {
	return &CurrencyMongo{db: db, i18n: i18n}
}

func (r *CurrencyMongo) CreateCurrency(userID string, data *domain.CurrencyInput) (domain.Currency, error) {
	var result domain.Currency

	collection := r.db.Collection(TblCurrency)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	// userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return nil, err
	// }
	// count, err := r.db.Collection(tblLanguage).CountDocuments(ctx, bson.M{})
	// if err != nil {
	// 	return response, err
	// }
	// newId := count + 1

	newPage := domain.Currency{
		Status:        data.Status,
		Title:         data.Title,
		Code:          data.Code,
		SymbolLeft:    data.SymbolLeft,
		SymbolRight:   data.SymbolRight,
		DecimalPlaces: data.DecimalPlaces,
		Value:         data.Value,
		SortOrder:     data.SortOrder,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	res, err := collection.InsertOne(ctx, newPage)
	if err != nil {
		return result, err
	}

	err = r.db.Collection(TblCurrency).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *CurrencyMongo) GetCurrency(id string) (domain.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result domain.Currency

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Currency{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(TblCurrency).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return domain.Currency{}, err
	}

	return result, nil
}

func (r *CurrencyMongo) FindCurrency(params domain.RequestParams) (domain.Response[domain.Currency], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Currency
	var response domain.Response[domain.Currency]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Currency]{}, err
	}
	cursor, err := r.db.Collection(TblCurrency).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Currency, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	var options options.CountOptions
	// options.SetLimit(params.Limit)
	options.SetSkip(params.Skip)
	count, err := r.db.Collection(TblCurrency).CountDocuments(ctx, params.Filter, &options)
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Currency]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *CurrencyMongo) UpdateCurrency(id string, data interface{}) (domain.Currency, error) {
	var result domain.Currency
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblCurrency)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": data})
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *CurrencyMongo) DeleteCurrency(id string) (domain.Currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = domain.Currency{}
	collection := r.db.Collection(TblCurrency)

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
