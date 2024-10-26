package domain

import (
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	// swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Login        string             `json:"login" bson:"login" binding:"required"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"password" bson:"password" binding:"required"`
	Strategy     string             `json:"strategy" bson:"strategy"`
	VkID         string             `json:"vkId" bson:"vk_id"`
	AppleID      string             `json:"appleId" bson:"apple_id"`
	GoogleID     string             `json:"googleId" bson:"google_id"`
	GithubID     string             `json:"githubId" bson:"github_id"`
	Verification Verification       `json:"verification" bson:"verification"`
	Session      Session            `json:"session" bson:"session"`

	Roles       []string `json:"roles" bson:"roles"`
	MaxDistance int      `json:"maxDistance" bson:"max_distance"`

	UserData  model.User `json:"-" bson:"user_data"`
	CreatedAt time.Time  `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" bson:"updated_at"`
}

type SignInInput struct {
	Login    string `json:"login" bson:"login" form:"login"`
	Email    string `json:"email" bson:"email" form:"email"`
	Password string `json:"password" bson:"password" form:"password"`
	Strategy string `json:"strategy" bson:"strategy"`
	VkID     string `json:"vkId" bson:"vk_id" form:"vkId"`
	AppleID  string `json:"appleId" bson:"apple_id" form:"appleId"`
	GoogleID string `json:"googleId" bson:"google_id" form:"googleId"`
	GithubID string `json:"githubId" bson:"github_id" form:"githubId"`

	Verification Verification `json:"verification" bson:"verification"`
	Session      Session      `json:"session" bson:"session"`
	Roles        []string     `json:"roles" bson:"roles"`
	MaxDistance  int          `json:"maxDistance" bson:"max_distance"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type DataForClaims struct {
	Roles  []string `json:"roles" bson:"roles"`
	UserID string   `json:"user_id" bson:"user_id"`
	Md     int      `json:"md" bson:"md"`
	UID    string   `json:"uid" bson:"uid"`
}

type Verification struct {
	Code     string `json:"code" bson:"code"`
	Verified bool   `json:"verified" bson:"verified"`
}
