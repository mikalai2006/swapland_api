package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"
	"fmt"

	"github.com/mikalai2006/swapland-api/graph/generated"
	"github.com/mikalai2006/swapland-api/graph/model"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tickets is the resolver for the tickets field.
func (r *queryResolver) Tickets(ctx context.Context, limit *int, skip *int, input *model.ParamsTicket) (*model.PaginationTicket, error) {
	var results *model.PaginationTicket

	filter := bson.D{}
	if input.UserID != nil {
		nID, _ := primitive.ObjectIDFromHex(*input.UserID)
		filter = append(filter, bson.E{"user_id", nID})
	}

	allItems, err := r.Repo.Ticket.FindTicket(domain.RequestParams{
		Options: domain.Options{Limit: int64(*limit), Sort: bson.D{{"updated_at", -1}}},
		Filter:  filter,
	})
	if err != nil {
		return results, err
	}

	data := make([]*model.Ticket, len(allItems.Data))
	for i, _ := range allItems.Data {
		data[i] = &allItems.Data[i]
	}

	total := len(data)

	results = &model.PaginationTicket{
		Data:  data,
		Total: &total,
		Limit: limit,
		Skip:  skip,
	}

	return results, nil
}

// Ticket is the resolver for the ticket field.
func (r *queryResolver) Ticket(ctx context.Context, input *model.ParamsTicket) (*model.Ticket, error) {
	var result *model.Ticket

	filter := bson.D{}
	if input.UserID != nil {
		nID, _ := primitive.ObjectIDFromHex(*input.UserID)
		filter = append(filter, bson.E{"_id", nID})
	}

	allItems, err := r.Repo.Ticket.FindTicket(domain.RequestParams{
		Options: domain.Options{Limit: 1},
		Filter:  filter,
	})
	if err != nil {
		return result, err
	}

	// data := make([]*model.Ticket, len(allItems.Data))
	// for i, _ := range allItems.Data {
	// 	data[i] = &allItems.Data[i]
	// }

	// total := len(data)

	result = &allItems.Data[0]

	return result, nil
}

// ID is the resolver for the id field.
func (r *ticketResolver) ID(ctx context.Context, obj *model.Ticket) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *ticketResolver) UserID(ctx context.Context, obj *model.Ticket) (string, error) {
	return obj.UserID.Hex(), nil
}

// ID is the resolver for the id field.
func (r *ticketMessageResolver) ID(ctx context.Context, obj *model.TicketMessage) (string, error) {
	return obj.ID.Hex(), nil
}

// UserID is the resolver for the userId field.
func (r *ticketMessageResolver) UserID(ctx context.Context, obj *model.TicketMessage) (string, error) {
	return obj.UserID.Hex(), nil
}

// TicketID is the resolver for the ticketId field.
func (r *ticketMessageResolver) TicketID(ctx context.Context, obj *model.TicketMessage) (string, error) {
	return obj.TicketID.Hex(), nil
}

// CreatedAt is the resolver for the createdAt field.
func (r *ticketMessageResolver) CreatedAt(ctx context.Context, obj *model.TicketMessage) (string, error) {
	return obj.CreatedAt.String(), nil
}

// UpdatedAt is the resolver for the updatedAt field.
func (r *ticketMessageResolver) UpdatedAt(ctx context.Context, obj *model.TicketMessage) (string, error) {
	return obj.UpdatedAt.String(), nil
}

// Ticket returns generated.TicketResolver implementation.
func (r *Resolver) Ticket() generated.TicketResolver { return &ticketResolver{r} }

// TicketMessage returns generated.TicketMessageResolver implementation.
func (r *Resolver) TicketMessage() generated.TicketMessageResolver { return &ticketMessageResolver{r} }

type ticketResolver struct{ *Resolver }
type ticketMessageResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *ticketResolver) Tilte(ctx context.Context, obj *model.Ticket) (string, error) {
	return obj.Title, nil
}
func (r *queryResolver) Tikets(ctx context.Context, limit *int, skip *int, input *model.ParamsTicket) (*model.PaginationTicket, error) {
	var results *model.PaginationTicket

	filter := bson.D{}
	if input.UserID != nil {
		nID, _ := primitive.ObjectIDFromHex(*input.UserID)
		filter = append(filter, bson.E{"node_id", nID})
	}

	allItems, err := r.Repo.Ticket.FindTicket(domain.RequestParams{
		Options: domain.Options{Limit: int64(*limit), Sort: bson.D{{"updated_at", -1}}},
		Filter:  filter,
	})
	if err != nil {
		return results, err
	}

	data := make([]*model.Ticket, len(allItems.Data))
	for i, _ := range allItems.Data {
		data[i] = &allItems.Data[i]
	}

	total := len(data)

	results = &model.PaginationTicket{
		Data:  data,
		Total: &total,
		Limit: limit,
		Skip:  skip,
	}

	return results, nil
}
func (r *queryResolver) Tiket(ctx context.Context, input *model.ParamsTicket) (*model.Ticket, error) {
	panic(fmt.Errorf("not implemented: Tiket - tiket"))
}
