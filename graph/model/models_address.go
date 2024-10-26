package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID string             `json:"userId" bson:"user_id"`

	OsmID    string                 `json:"osmId" bson:"osm_id" form:"osmId"`
	Address  map[string]interface{} `json:"address" bson:"address" form:"address"`
	DAddress string                 `json:"dAddress" bson:"d_address" form:"dAddress"`
	Lang     string                 `json:"lang" bson:"lang" form:"lang"`
	Props    map[string]interface{} `json:"props" bson:"props" form:"props"`

	// Lon       float64            `json:"lon" bson:"lon"`
	// Lat       float64            `json:"lat" bson:"lat"`
	// Type      string             `json:"type" bson:"type"`
	// Tags      interface{}        `json:"tags" bson:"tags"`
	// OsmID     string             `json:"osmId" bson:"osm_id"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
