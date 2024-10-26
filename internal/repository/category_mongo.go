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
)

type CategoryMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewCategoryMongo(db *mongo.Database, i18n config.I18nConfig) *CategoryMongo {
	return &CategoryMongo{db: db, i18n: i18n}
}

func (r *CategoryMongo) FindCategory(params domain.RequestParams) (domain.Response[model.Category], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Category
	var response domain.Response[model.Category]
	// filter, opts, err := CreateFilterAndOptions(params)
	// if err != nil {
	// 	return domain.Response[model.Category]{}, err
	// }
	// cursor, err := r.db.Collection(TblCategory).Find(ctx, filter, opts)
	pipe, err := CreatePipeline(params, &r.i18n)

	if err != nil {
		return response, err
	}

	cursor, err := r.db.Collection(TblCategory).Aggregate(ctx, pipe)
	// fmt.Println("filter Category:::", pipe)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Category, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblCategory).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Category]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *CategoryMongo) GetAllCategory(params domain.RequestParams) (domain.Response[model.Category], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []model.Category
	var response domain.Response[model.Category]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[model.Category]{}, err
	}

	cursor, err := r.db.Collection(TblCategory).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]model.Category, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(TblCategory).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[model.Category]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *CategoryMongo) CreateCategory(userID string, Category *model.Category) (*model.Category, error) {
	var result *model.Category
	collection := r.db.Collection(TblCategory)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	newCategory := model.Category{
		UserID:      userIDPrimitive,
		Seo:         Category.Seo,
		Title:       Category.Title,
		Description: Category.Description,
		Props:       Category.Props,
		Locale:      Category.Locale,
		Parent:      Category.Parent,
		Status:      Category.Status,
		SortOrder:   Category.SortOrder,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	res, err := collection.InsertOne(ctx, newCategory)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(TblCategory).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *CategoryMongo) GqlGetCategorys(params domain.RequestParams) ([]*model.Category, error) {
	fmt.Println("GqlGetCategorys")
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []*model.Category
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return results, err
	}
	// fmt.Println(pipe)

	cursor, err := r.db.Collection(TblCategory).Aggregate(ctx, pipe)
	if err != nil {
		return results, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return results, er
	}

	resultSlice := make([]*model.Category, len(results))

	copy(resultSlice, results)
	return results, nil
}

func (r *CategoryMongo) UpdateCategory(id string, userID string, data *model.CategoryInput) (*model.Category, error) {
	var result *model.Category
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	collection := r.db.Collection(TblCategory)

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, err
	}

	filter := bson.M{"_id": idPrimitive}
	// _, err = collection.UpdateOne(ctx, filter, bson.M{"$set": bson.M{
	// 	"seo":         data.Seo,
	// 	"title":       data.Title,
	// 	"description": data.Description,
	// 	"props":       data.Props,
	// 	"locale":      data.Locale,
	// 	"parent":      data.Parent,
	// 	"status":      data.Status,
	// 	"sort_order":  data.SortOrder,
	// 	"updated_at":  time.Now(),
	// }})
	// if err != nil {
	// 	return result, err
	// }
	newData := bson.M{}
	if data.Status != 0 {
		newData["status"] = data.Status
	}
	if data.SortOrder != 0 {
		newData["sort_order"] = data.SortOrder
	}
	if data.Title != "" {
		newData["title"] = data.Title
	}
	if data.Seo != "" {
		newData["seo"] = data.Seo
	}
	if data.Parent != "" {
		parentPrimitive, err := primitive.ObjectIDFromHex(data.Parent)
		if err != nil {
			return result, err
		}
		newData["parent"] = parentPrimitive
	}
	if data.Description != "" {
		newData["description"] = data.Description
	}
	if len(data.Locale) != 0 {
		newData["locale"] = data.Locale
	}
	if data.Props != nil {
		//newProps := make(map[string]interface{})
		newProps := data.Props
		if val, ok := data.Props["status"]; ok {
			if val == -1.0 {
				newDel := make(map[string]interface{})
				newDel["user_id"] = userID
				newDel["del_at"] = time.Now()
				newProps["del"] = newDel
			}
		}
		newData["props"] = newProps
	}
	newData["updated_at"] = time.Now()
	// test := model.ProductLike{}
	// if data.ProductLike != test {
	// 	newData["Product_like"] = data.ProductLike
	// }
	// if data.Status != 0 {
	// 	newData["status"] = data.Status
	// }
	// bson.M{
	// 	"lon":        data.Lon,
	// 	"lat":        data.Lat,
	// 	"type":       data.Type,
	// 	"osm_id":     data.OsmID,
	// 	"amenity_id": data.AmenityID,
	// 	"props":      data.Props,
	// 	"name":       data.Name,
	// 	"updated_at": time.Now(),
	// }
	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": newData})

	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *CategoryMongo) DeleteCategory(id string) (model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var result = model.Category{}
	collection := r.db.Collection(TblCategory)

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
