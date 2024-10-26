package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"

	"github.com/mikalai2006/swapland-api/graph/generated"
	"github.com/mikalai2006/swapland-api/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Address is the resolver for the address field.
func (r *geoLocationResolver) Address(ctx context.Context, obj *model.GeoLocation) (any, error) {
	return obj.Address, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id *string) (*model.User, error) {
	var result model.User
	// gc, err := utils.GinContextFromContext(ctx)
	// if err != nil {
	// 	return result, err
	// }
	// lang := gc.MustGet("i18nLocale").(string)

	filter := bson.D{}
	if id != nil {
		userIDPrimitive, err := primitive.ObjectIDFromHex(*id)
		if err != nil {
			return &result, err
		}

		filter = append(filter, bson.E{"_id", userIDPrimitive})
	}

	// allItems, err := r.Repo.User.GqlGetUsers(domain.RequestParams{
	// 	Options: domain.Options{Limit: 1, Skip: 0},
	// 	Filter:  filter,
	// 	Lang:    lang,
	// })
	// if err != nil {
	// 	return result, err
	// }

	// if len(allItems) > 0 {
	// 	result = allItems[0]
	// }
	result, err := r.Repo.User.GetUser(*id)
	if err != nil {
		return &result, err
	}

	return &result, nil
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *userResolver) UserID(ctx context.Context, obj *model.User) (string, error) {
	return obj.UserID.Hex(), nil
}

// GeoLocation returns generated.GeoLocationResolver implementation.
func (r *Resolver) GeoLocation() generated.GeoLocationResolver { return &geoLocationResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type geoLocationResolver struct{ *Resolver }
type userResolver struct{ *Resolver }