package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GeoLocation struct {
	Lon     float64                `json:"lon" bson:"lon"`
	Lat     float64                `json:"lat" bson:"lat"`
	OsmID   string                 `json:"osmId" bson:"osm_id"`
	Address map[string]interface{} `json:"address" bson:"address"`
}

type UserStat struct {
	AddProduct  int64 `json:"addProduct" bson:"addProduct"`
	TakeProduct int64 `json:"takeProduct" bson:"takeProduct"`
	GiveProduct int64 `json:"giveProduct" bson:"giveProduct"`
	AddOffer    int64 `json:"addOffer" bson:"addOffer"`
	TakeOffer   int64 `json:"takeOffer" bson:"takeOffer"`
	AddMessage  int64 `json:"addMessage" bson:"addMessage"`
	TakeMessage int64 `json:"takeMessage" bson:"takeMessage"`
	AddReview   int64 `json:"addReview" bson:"addReview"`
	TakeReview  int64 `json:"takeReview" bson:"takeReview"`

	Warning int64 `json:"warning" bson:"warning"`
	Request int64 `json:"request" bson:"request"`

	Subcribe    int64     `json:"subscribe" bson:"subscribe"`
	Subcriber   int64     `json:"subscriber" bson:"subscriber"`
	LastRequest time.Time `json:"lastRequest" bson:"lastRequest"`
}

type User struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty" primitive:"true"`
	UserID primitive.ObjectID `json:"userId,omitempty" bson:"user_id,omitempty" primitive:"true"`

	Name     string      `json:"name" bson:"name" form:"name"`
	Login    string      `json:"login" bson:"login" form:"login"`
	Currency string      `json:"currency" bson:"currency" form:"currency"`
	Lang     string      `json:"lang" bson:"lang" form:"lang"`
	Avatar   string      `json:"avatar" bson:"avatar"`
	Online   bool        `json:"online" bson:"online" form:"online"`
	Verify   bool        `json:"verify" bson:"verify"`
	Location GeoLocation `json:"location" bson:"location" form:"location"`

	UserStat UserStat `json:"userStat" bson:"user_stat"`
	// Test     interface{} `json:"test" bson:"test"`

	Roles  []string `json:"roles" bson:"roles"`
	Md     int      `json:"md" bson:"md"`
	Bal    int      `json:"bal" bson:"bal"`
	Images []Image  `json:"images,omitempty" bson:"images,omitempty"`

	LastTime  time.Time `json:"lastTime" bson:"last_time"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type UserInput struct {
	ID       string `json:"id" bson:"_id" form:"id" primitive:"true"`
	UserID   string `json:"userId" bson:"user_id" form:"userId" primitive:"true"`
	Name     string `json:"name" bson:"name" form:"name"`
	Login    string `json:"login" bson:"login" form:"login"`
	Currency string `json:"currency" bson:"currency" form:"currency"`
	Lang     string `json:"lang" bson:"lang" form:"lang"`
	Avatar   string `json:"avatar" bson:"avatar" form:"avatar"`
}
