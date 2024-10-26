package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NodeVote struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodeUserID primitive.ObjectID `json:"nodeUserId" bson:"node_user_id"`
	NodeID     primitive.ObjectID `json:"nodeId" bson:"node_id"`
	Value      int                `json:"value" bson:"value"`
	User       User               `json:"user,omitempty" bson:"user,omitempty"`
	Owner      User               `json:"owner,omitempty" bson:"owner,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

// type NodeVoteMongo struct {
// 	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
// 	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

// 	NodedataID primitive.ObjectID `json:"nodedataId" bson:"nodedata_id"`
// 	Value      int                `json:"value" bson:"value"`
// 	// Status     int64              `json:"status" bson:"status"`

// 	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
// 	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
// }

type NodeVoteInput struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodeUserID primitive.ObjectID `json:"nodeUserId" bson:"node_user_id"`
	NodeID     primitive.ObjectID `json:"nodeId" bson:"node_id"`
	Value      int                `json:"value" bson:"value"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
