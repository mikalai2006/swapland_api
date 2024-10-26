package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Library struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id" primitive:"true"`
	Name   string             `json:"name" bson:"name"`
	Icon   string             `json:"icon" bson:"icon"`
	Groups string             `json:"groups" bson:"groups"`

	Tree      interface{} `json:"tree" bson:"tree"`
	Data      []*Field    `json:"-" bson:"data"`
	CreatedAt time.Time   `json:"createdAt" bson:"createdAt" form:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt" bson:"updatedAt" form:"updatedAt"`
}

type LibraryInput struct {
	Name   string `json:"name" bson:"name" form:"name"`
	Icon   string `json:"icon" bson:"icon" form:"icon"`
	Groups string `json:"groups" bson:"groups" form:"groups"`

	Data interface{} `json:"tree" bson:"tree" form:"data"`
}

type Field struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id" primitive:"true"`
	UID       string             `json:"_uid" bson:"_uid"`
	Parent    string             `json:"parent" bson:"parent"`
	Publish   bool               `json:"publish" bson:"publish"`
	Name      string             `json:"name" bson:"name"`
	UserID    primitive.ObjectID `json:"user_id" bson:"user_id"  primitive:"true"`
	LibraryID primitive.ObjectID `json:"libraryId" bson:"libraryId"  primitive:"true"`

	Data      FieldNode `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type FieldValue struct {
	UID    string        `json:"_uid" bson:"_uid"`
	Parent string        `json:"parent" bson:"parent"`
	Name   string        `json:"name" bson:"name"`
	Child  []interface{} `json:"child" bson:"child"`
}

type FieldNode struct {
	Value map[string]any `json:"value" bson:"value"`
}
