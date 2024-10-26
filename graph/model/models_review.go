package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaginationReview struct {
	Total int       `json:"total,omitempty"`
	Limit int       `json:"limit,omitempty"`
	Skip  int       `json:"skip,omitempty"`
	Data  []*Review `json:"data,omitempty"`
}

type Review struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	NodeID primitive.ObjectID `json:"nodeId" bson:"node_id"`

	Review string `json:"review" bson:"review"`
	Rate   int    `json:"rate" bson:"rate"`

	Publish   bool      `json:"publish" bson:"publish"`
	User      User      `json:"user,omitempty" bson:"user,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ReviewInput struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`
	NodeID primitive.ObjectID `json:"nodeId" bson:"node_id"`

	Review string `json:"review" bson:"review"`
	Rate   int    `json:"rate" bson:"rate"`

	Publish   bool      `json:"publish" bson:"publish"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type ReviewInputData struct {
	NodeID string `json:"nodeId" bson:"node_id" form:"osmId"`

	Review string `json:"review" bson:"review"  form:"review"`
	Rate   int    `json:"rate" bson:"rate" form:"rate"`

	Publish   bool      `json:"publish" bson:"publish" form:"publish"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
