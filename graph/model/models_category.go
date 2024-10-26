package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      primitive.ObjectID           `json:"userId" bson:"user_id"`
	Seo         string                       `json:"seo" bson:"seo"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Parent      primitive.ObjectID           `json:"parent" bson:"parent"`
	Status      int64                        `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	SortOrder   int64                        `json:"sortOrder" bson:"sort_order"`
	CreatedAt   time.Time                    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                    `json:"updatedAt" bson:"updated_at"`
}

type CategoryInput struct {
	ID          primitive.ObjectID           `json:"id" bson:"_id,omitempty"`
	UserID      string                       `json:"userId" bson:"user_id" form:"userId"`
	Seo         string                       `json:"seo" bson:"seo" form:"seo"`
	Title       string                       `json:"title" bson:"title" form:"title"`
	Description string                       `json:"description" bson:"description" form:"description"`
	Props       map[string]interface{}       `json:"props" bson:"props"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Parent      string                       `json:"parent" bson:"parent" form:"parent"`
	Status      int64                        `json:"status" bson:"status" form:"status"`
	SortOrder   int64                        `json:"sortOrder" bson:"sort_order" form:"sortOrder"`
	CreatedAt   time.Time                    `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time                    `json:"updatedAt" bson:"updated_at"`
}
