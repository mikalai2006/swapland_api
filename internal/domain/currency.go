package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Currency struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Status        bool               `json:"status" bson:"status"`
	Title         string             `json:"title" bson:"title"`
	Code          string             `json:"code" bson:"code"`
	SymbolLeft    string             `json:"symbolLeft" bson:"symbol_left"`
	SymbolRight   string             `json:"symbolRight" bson:"symbol_right"`
	DecimalPlaces int64              `json:"decimalPlaces" bson:"decimal_places"`
	Value         float64            `json:"value" bson:"value"`
	SortOrder     int64              `json:"sortOrder" bson:"sort_order"`
	CreatedAt     time.Time          `json:"createdAt" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updated_at"`
}

type CurrencyInput struct {
	Status        bool    `json:"status" bson:"status"`
	Title         string  `json:"title" bson:"title"`
	Code          string  `json:"code" bson:"code"`
	SymbolLeft    string  `json:"symbolLeft" bson:"symbol_left"`
	SymbolRight   string  `json:"symbolRight" bson:"symbol_right"`
	DecimalPlaces int64   `json:"decimalPlaces" bson:"decimal_places"`
	Value         float64 `json:"value" bson:"value"`
	SortOrder     int64   `json:"sortOrder" bson:"sort_order"`
}
