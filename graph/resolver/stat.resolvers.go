package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"
	"fmt"

	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StatNode is the resolver for the statNode field.
func (r *queryResolver) StatNode(ctx context.Context) (*model.StatNode, error) {
	var result *model.StatNode

	pipe := mongo.Pipeline{}

	pipe = append(pipe, bson.D{{"$group", bson.M{"_id": bson.D{{"$toString", "$ccode"}}, "count": bson.M{"$sum": 1}}}})

	var allItems []model.GroupNodeCountry
	cursor, err := r.DB.Collection(repository.TblNode).Aggregate(ctx, pipe)
	if err != nil {
		return result, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &allItems); er != nil {
		return result, er
	}

	pipe2 := mongo.Pipeline{}

	pipe2 = append(pipe2, bson.D{{"$group", bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}}})

	var nodeType []model.GroupNodeType
	cursor2, err := r.DB.Collection(repository.TblNode).Aggregate(ctx, pipe2)
	if err != nil {
		return result, err
	}
	defer cursor2.Close(ctx)

	if er := cursor2.All(ctx, &nodeType); er != nil {
		return result, er
	}

	fmt.Println(allItems)
	result = &model.StatNode{
		GroupCountry: allItems,
		GroupType:    nodeType,
	}

	return result, nil
}