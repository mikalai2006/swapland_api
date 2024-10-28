package repository

import (
	"context"
	"fmt"
	"strconv"
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

func (r *AddressMongo) FindAddress(input *model.AddressFilter) (domain.Response[model.Address], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Address
	var response domain.Response[model.Address]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[domain.Address]{}, err
	// }
	q := bson.D{}
	if input.OsmID != nil {
		q = append(q, bson.E{"osmId", *input.OsmID})
	}
	if input.Lat != 0 {
		q = append(q, bson.E{"lat", input.Lat})
	}
	if input.Lon != 0 {
		q = append(q, bson.E{"lon", input.Lon})
	}
	if input.UserID != nil {
		userIDPrimitive, err := primitive.ObjectIDFromHex(*input.UserID)
		if err != nil {
			return response, err
		}
		q = append(q, bson.E{"userId", userIDPrimitive})
	}

	pipe := mongo.Pipeline{}
	pipe = append(pipe, bson.D{{"$match", q}})

	limit := 100
	skip := 0
	if input.Limit != nil {
		limit = *input.Limit
	}
	if input.Skip != nil {
		skip = *input.Skip
	}

	pipe = append(pipe, bson.D{{"$limit", skip + limit}})
	pipe = append(pipe, bson.D{{"$skip", skip}})

	cursor, err := r.db.Collection(TblAddress).Aggregate(ctx, pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Address, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblAddress).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Address]{
		Total: int(count),
		Skip:  skip,
		Limit: limit,
		Data:  resultSlice,
	}
	return response, nil
}

func (r *AddressMongo) CreateAddress(userID string, address domain.ResponseNominatim) (*model.Address, error) {
	var result *model.Address

	collection := r.db.Collection(TblAddress)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	var lat float64
	if s, err := strconv.ParseFloat(address.Lat, 32); err == nil {
		lat = s
	}

	var lon float64
	if s, err := strconv.ParseFloat(address.Lon, 32); err == nil {
		lon = s
	}
	newAddress := model.Address{
		UserID:   userIDPrimitive,
		OsmID:    fmt.Sprintf("%v/%v", address.OsmType, address.OsmID),
		Lat:      lat,
		Lon:      lon,
		Lang:     address.Lang,
		Address:  address.Address,
		DAddress: address.DisplayName,
		//Props:     address.Address,
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

func (r *AddressMongo) UpdateAddress(id string, userID string, data domain.ResponseNominatim) (*model.Address, error) {
	var result *model.Address
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblAddress)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	// idUser, err := primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	return result, err
	// }
	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	// if data.Lat != 0 {
	// 	newData["lat"] = data.Lat
	// }
	// if data.Lon != 0 {
	// 	newData["lon"] = data.Lon
	// }
	// if data.OsmID != "" {
	// 	newData["osmId"] = data.OsmID
	// }
	// if data.Address != nil {
	// 	newData["address"] = data.Address
	// }
	// if data.DAddress != "" {
	// 	newData["dAddress"] = data.DAddress
	// }

	if s, err := strconv.ParseFloat(data.Lat, 32); err == nil {
		newData["lat"] = s
	}

	if s, err := strconv.ParseFloat(data.Lon, 32); err == nil {
		newData["lon"] = s
	}

	if data.OsmID != 0 {
		newData["osmId"] = fmt.Sprintf("%v/%v", data.OsmType, data.OsmID)
	}

	if data.DisplayName != "" {
		newData["dAddress"] = data.DisplayName
	}
	if data.Address != nil {
		newData["address"] = data.Address
	}

	// if data.Props != nil {
	// 	//newProps := make(map[string]interface{})
	// 	newProps := data.Props
	// 	if val, ok := data.Props["status"]; ok {
	// 		if val == -1.0 {
	// 			newDel := make(map[string]interface{})
	// 			newDel["user_id"] = userID
	// 			newDel["del_at"] = time.Now()
	// 			newProps["del"] = newDel
	// 		}
	// 	}
	// 	newData["props"] = newProps
	// }
	newData["updatedAt"] = time.Now()

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
