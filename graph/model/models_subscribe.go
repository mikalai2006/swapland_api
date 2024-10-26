package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subscribe struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userId" bson:"user_id"`
	SubUserID primitive.ObjectID `json:"subUserId" bson:"sub_user_id"`
	Status    int                `json:"status" bson:"status"`

	User    User `json:"user" bson:"user,omitempty"`
	SubUser User `json:"subUser" bson:"sub_user,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type SubscribeInput struct {
	UserID    primitive.ObjectID `json:"userId" bson:"user_id"`
	SubUserID primitive.ObjectID `json:"subUserId" bson:"sub_user_id"`
	Status    int                `json:"status" bson:"status" form:"status"`
	CreatedAt time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updated_at"`
}
