package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"github.com/mikalai2006/swapland-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewUserMongo(db *mongo.Database, i18n config.I18nConfig) *UserMongo {
	return &UserMongo{db: db, i18n: i18n}
}

func (r *UserMongo) Iam(userID string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result model.User
	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return model.User{}, err
	}

	params := domain.RequestParams{}
	params.Filter = bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, params.Filter).Decode(&result)
	if err != nil {
		return model.User{}, err
	}

	pipe, err := CreatePipeline(params, &r.i18n) // mongo.Pipeline{bson.D{{"_id", userIDPrimitive}}} //
	if err != nil {
		return result, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
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
	}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // .FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(&result); er != nil {
			return result, er
		}
	}

	return result, nil
}

func (r *UserMongo) GetUser(id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result model.User

	userIDPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = r.db.Collection(tblUsers).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return model.User{}, err
	}

	pipe, err := CreatePipeline(domain.RequestParams{
		Filter: filter,
	}, &r.i18n)
	if err != nil {
		return result, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
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
	}})
	// add populate.
	pipe = append(pipe, bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from": TblAuth,
			"as":   "authsx",
			// "localField":   "_id",
			// "foreignField": "service_id",
			"let": bson.D{{Key: "userId", Value: "$user_id"}},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$_id", "$$userId"}}}}},
				bson.D{{"$limit", 1}},
			},
		},
	}})

	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.M{"$first": "$authsx"}}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"md": "$test.max_distance"}}})
	pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"roles": "$test.roles"}}})

	// // add stat user tag vote.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodedataVote,
	// 		"as":   "tests",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 			bson.D{
	// 				{
	// 					"$group", bson.D{
	// 						{
	// 							"_id", "",
	// 						},
	// 						{"valueTagLike", bson.D{{"$sum", "$value"}}},
	// 						{"countTagLike", bson.D{{"$sum", 1}}},
	// 					},
	// 				},
	// 			},
	// 			bson.D{{Key: "$project", Value: bson.M{"_id": 0, "valueTagLike": "$valueTagLike", "countTagLike": "$countTagLike"}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.M{"$first": "$tests"}}}})

	// // add stat user node votes.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodeVote,
	// 		"as":   "tests2",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 			bson.D{
	// 				{
	// 					"$group", bson.D{
	// 						{
	// 							"_id", "",
	// 						},
	// 						{"valueNodeLike", bson.D{{"$sum", "$value"}}},
	// 						{"countNodeLike", bson.D{{"$sum", 1}}},
	// 					},
	// 				},
	// 			},
	// 			bson.D{{Key: "$project", Value: bson.M{"_id": 0, "valueNodeLike": "$valueNodeLike", "countNodeLike": "$countNodeLike"}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"$first": "$tests2"},
	// 	},
	// }},
	// }}})

	// // add stat user node votes.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNode,
	// 		"as":   "countNodes",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"countNodes": bson.M{"$size": "$countNodes"}},
	// 	},
	// }},
	// }}})

	// // add stat user added nodedata.
	// pipe = append(pipe, bson.D{{
	// 	Key: "$lookup",
	// 	Value: bson.M{
	// 		"from": TblNodedata,
	// 		"as":   "countNodedatas",
	// 		"let":  bson.D{{Key: "userId", Value: "$_id"}},
	// 		"pipeline": mongo.Pipeline{
	// 			bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{"$eq": [2]string{"$user_id", "$$userId"}}}}},
	// 		},
	// 	},
	// }})
	// pipe = append(pipe, bson.D{{Key: "$set", Value: bson.M{"test": bson.D{{
	// 	"$mergeObjects", bson.A{
	// 		"$test",
	// 		bson.M{"countNodedatas": bson.M{"$size": "$countNodedatas"}},
	// 	},
	// }},
	// }}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if er := cursor.Decode(&result); er != nil {
			return result, er
		}
	}

	return result, nil
}

func (r *UserMongo) FindUser(params domain.RequestParams) (domain.Response[model.User], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.User
	var response domain.Response[model.User]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.User]{}, err
	}
	fmt.Println("params:::", params)

	// add populate.
	pipe = append(pipe, bson.D{{
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
	}})

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.User, len(results))
	copy(resultSlice, results)

	count, err := r.db.Collection(tblUsers).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.User]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *UserMongo) CreateUser(userID string, user *model.User) (*model.User, error) {
	var result *model.User

	collection := r.db.Collection(tblUsers)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newUser := model.User{
		Avatar:    user.Avatar,
		Name:      user.Name,
		UserID:    userIDPrimitive,
		Login:     user.Login,
		Lang:      user.Lang,
		Currency:  user.Currency,
		Online:    user.Online,
		Verify:    user.Verify,
		Bal:       3,
		LastTime:  time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblUsers).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *UserMongo) DeleteUser(id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.User{}
	collection := r.db.Collection(tblUsers)

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

func (r *UserMongo) UpdateUser(id string, user *model.User) (model.User, error) {
	var result model.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	// data, err := utils.GetBodyToData(user)
	// if err != nil {
	// 	return result, err
	// }

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}

	newData := bson.M{}
	if user.Lang != "" {
		newData["lang"] = user.Lang
	}
	if user.Name != "" {
		newData["name"] = user.Name
	}
	if user.Login != "" {
		newData["login"] = user.Login
	}
	if user.Location.Lat != 0 {
		newData["location"] = user.Location
	}

	newData["online"] = user.Online

	// fmt.Println("data=", user)
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})
	if err != nil {
		return result, err
	}

	// err = collection.FindOne(ctx, filter).Decode(&result)
	results, err := r.FindUser(domain.RequestParams{Filter: bson.M{"_id": idPrimitive}, Options: domain.Options{Limit: 1}})
	if err != nil {
		return result, err
	}
	result = results.Data[0]

	return result, nil
}

func (r *UserMongo) GqlGetUsers(params domain.RequestParams) ([]*model.User, error) {
	fmt.Println("GqlGetUsers")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.User
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}

	// add populate.
	pipe = append(pipe, bson.D{{
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
	}})

	// fmt.Println(pipe)

	cursor, err := r.db.Collection(tblUsers).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.User, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *UserMongo) SetStat(userID string, inputData model.UserStat) (model.User, error) {
	var result model.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	newData := bson.M{}
	if inputData.AddProduct != 0 {
		newData["user_stat.addProduct"] = utils.Max(result.UserStat.AddProduct+inputData.AddProduct, 0)
	}
	if inputData.TakeProduct != 0 {
		newData["user_stat.takeProduct"] = utils.Max(result.UserStat.TakeProduct+inputData.TakeProduct, 0)
	}
	if inputData.GiveProduct != 0 {
		newData["user_stat.giveProduct"] = utils.Max(result.UserStat.GiveProduct+inputData.GiveProduct, 0)
	}
	if inputData.AddOffer != 0 {
		newData["user_stat.addOffer"] = utils.Max(result.UserStat.AddOffer+inputData.AddOffer, 0)
	}
	if inputData.TakeOffer != 0 {
		newData["user_stat.takeOffer"] = utils.Max(result.UserStat.TakeOffer+inputData.TakeOffer, 0)
	}
	if inputData.AddMessage != 0 {
		newData["user_stat.addMessage"] = utils.Max(result.UserStat.AddMessage+inputData.AddMessage, 0)
	}
	if inputData.TakeMessage != 0 {
		newData["user_stat.takeMessage"] = utils.Max(result.UserStat.TakeMessage+inputData.TakeMessage, 0)
	}
	if inputData.AddReview != 0 {
		newData["user_stat.addReview"] = utils.Max(result.UserStat.AddReview+inputData.AddReview, 0)
	}
	if inputData.TakeReview != 0 {
		newData["user_stat.takeReview"] = utils.Max(result.UserStat.TakeReview+inputData.TakeReview, 0)
	}
	if inputData.Warning != 0 {
		newData["user_stat.warning"] = utils.Max(result.UserStat.Warning+inputData.Warning, 0)
	}
	if inputData.Request != 0 {
		newData["user_stat.request"] = utils.Max(result.UserStat.Request+inputData.Request, 0)
	}
	if inputData.Subcribe != 0 {
		newData["user_stat.subscribe"] = utils.Max(result.UserStat.Subcribe+inputData.Subcribe, 0)
	}
	if inputData.Subcriber != 0 {
		newData["user_stat.subscriber"] = utils.Max(result.UserStat.Subcriber+inputData.Subcriber, 0)
	}
	if !inputData.LastRequest.IsZero() {
		newData["user_stat.lastRequest"] = result.UserStat.LastRequest
	}

	// fmt.Println("newData=", newData)
	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": newData}).Decode(&result)
	if err != nil {
		return result, err
	}

	// var operations []mongo.WriteModel
	// operationA := mongo.NewUpdateOneModel()
	// operationA.SetFilter(bson.M{"_id": userIDPrimitive})
	// operationA.SetUpdate(bson.D{
	// 	{"$inc", bson.D{
	// 		{"user_stat.node", 1},
	// 	}},
	// })
	// operations = append(operations, operationA)
	// _, err = r.db.Collection(TblNode).BulkWrite(ctx, operations)

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *UserMongo) SetBal(userID string, value int) (model.User, error) {
	var result model.User
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(tblUsers)

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": userIDPrimitive}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	newData := bson.M{}
	if value != 0 {
		newData["bal"] = (int64)(result.Bal + value) //utils.Max((int64)(result.Bal+value), 0)
	}

	// fmt.Println("newData=", newData)
	err = collection.FindOneAndUpdate(ctx, filter, bson.M{"$set": newData}).Decode(&result)
	if err != nil {
		return result, err
	}

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
