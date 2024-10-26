package domain

import "time"

type Response[D any] struct {
	Total int `json:"total" bson:"total"`
	Limit int `json:"limit" bson:"limit"`
	Skip  int `json:"skip" bson:"skip"`
	Data  []D `json:"data" bson:"data"`
}

type ResponseTokens struct {
	AccessToken  string `json:"access_token" bson:"access_token"`
	RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in" bson:"expires_in"`
	ExpiresInR   int64  `json:"expires_in_r" bson:"expires_in_r"`
}

type GeneralFieldDB struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at" form:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" form:"updated_at"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseNominatim struct {
	OsmType     string                 `json:"osm_type" bson:"osm_type"`
	OsmID       int                    `json:"osm_id" bson:"osm_id"`
	Lat         string                 `json:"lat" bson:"lat"`
	Lon         string                 `json:"lon" bson:"lon"`
	Class       string                 `json:"class" bson:"class"`
	Type        string                 `json:"type" bson:"type"`
	AddressType string                 `json:"address_type" bson:"address_type"`
	Name        string                 `json:"name" bson:"name"`
	DisplayName string                 `json:"display_name" bson:"display_name"`
	Address     map[string]interface{} `json:"address" bson:"address"`
}

// type ResponseNominatimAddress struct {
// 	Highway      string `json:"highway"`
// 	Road         string `json:"road"`
// 	Village      string `json:"village"`
// 	Municipality string `json:"municipality"`
// 	County       string `json:"county"`
// 	State        string `json:"state"`
// 	Postcode     string `json:"postcode"`
// 	Country      string `json:"country"`
// 	CountryCode  string `json:"country_code"`
// }
