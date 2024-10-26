package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id" primitive:"true"`

	OsmID    string                 `json:"osmId" bson:"osm_id" form:"osmId"`
	Address  map[string]interface{} `json:"address" bson:"address" form:"address"`
	DAddress string                 `json:"dAddress" bson:"d_address" form:"dAddress"`
	Lang     string                 `json:"lang" bson:"lang" form:"lang"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at" form:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at" form:"updatedAt"`
}

type AddressInput struct {
	OsmID    string                 `json:"osmId" bson:"osm_id" form:"osmId" binding:"required"`
	Address  map[string]interface{} `json:"address" bson:"address" form:"address" binding:"required"`
	DAddress string                 `json:"dAddress" bson:"d_address" form:"dAddress" binding:"required"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`
	Lang     string                 `json:"lang" bson:"lang"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
