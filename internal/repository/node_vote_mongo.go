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

type NodeVoteMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewNodeVoteMongo(db *mongo.Database, i18n config.I18nConfig) *NodeVoteMongo {
	return &NodeVoteMongo{db: db, i18n: i18n}
}

func (r *NodeVoteMongo) FindNodeVote(params domain.RequestParams) (domain.Response[model.NodeVote], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.NodeVote
	var response domain.Response[model.NodeVote]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.NodeVote]{}, err
	// }
	// cursor, err := r.db.Collection(TblNodeVote).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
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

	// get owner node.
	pipe = append(pipe, bson.D{{Key: "$lookup", Value: bson.M{
		"from": "users",
		"as":   "usero",
		"let":  bson.D{{Key: "nodeUserId", Value: "$node_user_id"}},
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$nodeUserId"}}}}},
			bson.D{{"$limit", 1}},
			bson.D{{
				Key: "$lookup",
				Value: bson.M{
					"from": "image",
					"as":   "images",
					"let":  bson.D{{Key: "serviceId", Value: bson.D{{"$toString", "$_id"}}}},
					"pipeline": mongo.Pipeline{
						bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$service_id", "$$serviceId"}}}}},
					},
				},
			}},
		},
	}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"owner": bson.M{"$first": "$usero"}}}})

	cursor, err := r.db.Collection(TblNodeVote).Aggregate(ctx, pipe)
	// fmt.Println("filter NodeVote:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.NodeVote, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblNodeVote).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.NodeVote]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *NodeVoteMongo) CreateNodeVote(userID string, data *model.NodeVote) (*model.NodeVote, error) {
	var result *model.NodeVote

	collection := r.db.Collection(TblNodeVote)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	// nodedataIDPrimitive, err := primitive.ObjectIDFromHex(data.NodedataID)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println(bson.M{"nodedata_id": nodedataIDPrimitive, "user_id": userIDPrimitive})

	var existVote model.NodeVoteInput
	r.db.Collection(TblNodeVote).FindOne(ctx, bson.M{"node_id": data.NodeID, "user_id": userIDPrimitive}).Decode(&existVote)
	// if err != nil {
	// 	if errors.Is(err, mongo.ErrNoDocuments) {
	// 		return result, model.ErrAddressNotFound
	// 	}
	// 	return nil, err
	// }

	if (existVote == model.NodeVoteInput{}) {
		newNodeVote := model.NodeVoteInput{
			UserID:     userIDPrimitive,
			NodeID:     data.NodeID,
			Value:      data.Value,
			NodeUserID: data.NodeUserID,
			// Status:     100, //data.Status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		res, err := collection.InsertOne(ctx, newNodeVote)
		if err != nil {
			return nil, err
		}

		// err = r.db.Collection(TblNodeVote).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
		items, err := r.FindNodeVote(domain.RequestParams{
			Filter: bson.D{
				{"_id", res.InsertedID},
			},
			Options: domain.Options{
				Limit: 1,
			},
		})
		if err != nil {
			return nil, err
		}
		result = &items.Data[0]

		// statFragment := bson.D{}
		// if result.Value > 0 {
		// 	statFragment = append(statFragment, bson.E{"user_stat.poiLike", 1})
		// } else {
		// 	statFragment = append(statFragment, bson.E{"user_stat.poiDLike", 1})
		// }
		// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
		// 	{"$inc", statFragment},
		// })
	} else {
		updateNodeVote := &model.NodeVoteInput{
			// UserID:     userIDPrimitive,
			// NodedataID: nodedataIDPrimitive,
			Value: data.Value,
		}
		result, err = r.UpdateNodeVote(existVote.ID.Hex(), userID, updateNodeVote)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (r *NodeVoteMongo) UpdateNodeVote(id string, userID string, data *model.NodeVoteInput) (*model.NodeVote, error) {
	var result *model.NodeVote
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblNodeVote)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if data.Value != 0 {
		newData["value"] = data.Value
	}
	// if data.Status != 0 {
	// 	newData["status"] = data.Status
	// }
	newData["updated_at"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	// if err != nil {
	// 	return result, err
	// }
	items, err := r.FindNodeVote(domain.RequestParams{
		Filter: bson.D{
			{"_id", idPrimitive},
		},
		Options: domain.Options{
			Limit: 1,
		},
	})
	if err != nil {
		return nil, err
	}
	result = &items.Data[0]

	return result, nil
}

func (r *NodeVoteMongo) DeleteNodeVote(id string) (model.NodeVote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.NodeVote{}
	collection := r.db.Collection(TblNodeVote)

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

	// statFragment := bson.D{}
	// if result.Value > 0 {
	// 	statFragment = append(statFragment, bson.E{"user_stat.poiLike", -1})
	// } else {
	// 	statFragment = append(statFragment, bson.E{"user_stat.poiDLike", -1})
	// }
	// _, _ = r.db.Collection(tblUsers).UpdateOne(ctx, bson.M{"_id": result.UserID}, bson.D{
	// 	{"$inc", statFragment},
	// })

	return result, nil
}

// func (r *NodeVoteMongo) GetAllNodeVote(params domain.RequestParams) (domain.Response[model.NodeVote], error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
// 	defer cancel()

// 	var results []model.NodeVote
// 	var response domain.Response[model.NodeVote]
// 	pipe, err := CreatePipeline(params, &r.i18n)
// 	if err != nil {
// 		return domain.Response[model.NodeVote]{}, err
// 	}

// 	cursor, err := r.db.Collection(TblNodeVote).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
// 	if err != nil {
// 		return response, err
// 	}
// 	defer cursor.Close(ctx)

// 	if er := cursor.All(ctx, &results); er != nil {
// 		return response, er
// 	}

// 	resultSlice := make([]model.NodeVote, len(results))
// 	// for i, d := range results {
// 	// 	resultSlice[i] = d
// 	// }
// 	copy(resultSlice, results)

// 	count, err := r.db.Collection(TblNodeVote).CountDocuments(ctx, bson.M{})
// 	if err != nil {
// 		return response, err
// 	}

// 	response = domain.Response[model.NodeVote]{
// 		Total: int(count),
// 		Skip:  int(params.Options.Skip),
// 		Limit: int(params.Options.Limit),
// 		Data:  resultSlice,
// 	}
// 	return response, nil
// }
