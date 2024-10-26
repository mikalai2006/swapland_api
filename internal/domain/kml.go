package domain

import (
	"encoding/xml"
	"time"

	"github.com/mikalai2006/swapland-api/graph/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KMLParseSchema struct {
	XMLName xml.Name `xml:"kml"`
	// Text     string   `xml:",chardata"`
	XMLNS    string `xml:"xmlns,attr"`
	GX       string `xml:"xmlns:gx,attr"`
	KML      string `xml:"xmlns:kml,attr"`
	Atom     string `xml:"xmlns:atom,attr"`
	Document struct {
		ID            string        `xml:"id,attr"`
		KMLFileSchema KMLFileSchema `xml:"Schema"`
		KMLGX         []KMLGX       `xml:"gx:CascadingStyle"`
		StyleMap      []StyleMap    `xml:"StyleMap"`
		// Text string `xml:",chardata"`
		Name string `xml:"name"`
		// Snippet string `xml:"Snippet"`
		// Open   string      `xml:"open"`
		Folder []KMLFolder `xml:"Folder"`
	} `xml:"Document"`
}

type KMLSchema struct {
	XMLName xml.Name `xml:"kml"`
	// Text     string   `xml:",chardata"`
	XMLNS    string `xml:"xmlns,attr"`
	GX       string `xml:"xmlns:gx,attr"`
	KML      string `xml:"xmlns:kml,attr"`
	Atom     string `xml:"xmlns:atom,attr"`
	Document struct {
		ID            string        `xml:"id,attr"`
		KMLFileSchema KMLFileSchema `xml:"Schema"`
		KMLGX         []KMLGX       `xml:"gx:CascadingStyle"`
		StyleMap      []StyleMap    `xml:"StyleMap"`
		// Text string `xml:",chardata"`
		Name string `xml:"name"`
		// Snippet string `xml:"Snippet"`
		// Open   string      `xml:"open"`
		Folder KMLParentFolder `xml:"Folder"`
	} `xml:"Document"`
}

type KMLIgoicon struct {
	Filename string `xml:"filename"`
}

type KMLMetadata struct {
	Igoicon KMLIgoicon `xml:"igoicon"`
}

type KMLPoint struct {
	Coordinates string `xml:"coordinates"`
}

type KMLGX struct {
	ID       string     `xml:"kml:id,attr"`
	StyleURL string     `xml:"styleURL"`
	Style    KMLGXStyle `xml:"Style"`
}
type KMLGXStyle struct {
	IconStyle  KMLGXIconStyle `xml:"IconStyle"`
	LabelStyle KMLLabelStyle  `xml:"LabelStyle"`
	LineStyle  KMLLineStyle   `xml:"LineStyle"`
	PolyStyle  KMLPolyStyle   `xml:"PolyStyle"`
	// BalloonStyle KMLBalloonStyle `xml:"BalloonStyle"`
}
type KMLBalloonStyle struct {
	// Text        string `xml:"text"`
	DisplayMode string `xml:"gx:displayMode"`
}
type KMLGXIconStyle struct {
	Icon    KMLGXIcon `xml:"Icon"`
	Scale   string    `xml:"scale"`
	HotSpot HotSpot   `xml:"hotSpot"`
}
type KMLGXIcon struct {
	Href string `xml:"href"`
}
type HotSpot struct {
	X      string `xml:"x,attr"`
	Y      string `xml:"y,attr"`
	XUnits string `xml:"xunits,attr"`
	YUnits string `xml:"yunits,attr"`
}
type KMLLabelStyle struct {
	Scale string `xml:"scale"`
}
type KMLLineStyle struct {
	Color string `xml:"color"`
	Width string `xml:"width"`
}
type KMLPolyStyle struct {
	Color string `xml:"color"`
}
type StyleMap struct {
	ID   string    `xml:"id,attr"`
	Pair []KMLPair `xml:"Pair"`
}
type KMLPair struct {
	KEY      string `xml:"key"`
	StyleURL string `xml:"styleUrl"`
}
type KMLPlacemarkDescription struct {
	Description string `xml:",innerxml"`
}
type KMLPlacemark struct {
	ID           string                  `xml:"id,attr"`
	Name         string                  `xml:"name"`
	Description  KMLPlacemarkDescription `xml:"description"`
	StyleURL     string                  `xml:"styleUrl"`
	ExtendedData ExtendedData            `xml:"ExtendedData"`
	Point        KMLPoint                `xml:"Point"`
	// Style struct {
	// 	Text      string `xml:",chardata"`
	// 	LineStyle struct {
	// 		Text  string `xml:",chardata"`
	// 		Color string `xml:"color"`
	// 		Width string `xml:"width"`
	// 	} `xml:"LineStyle"`
	// 	PolyStyle struct {
	// 		Text    string `xml:",chardata"`
	// 		Color   string `xml:"color"`
	// 		Fill    string `xml:"fill"`
	// 		Outline string `xml:"outline"`
	// 	} `xml:"PolyStyle"`
	// } `xml:"Style"`
	// LineString struct {
	// 	Text        string `xml:",chardata"`
	// 	Tessellate  string `xml:"tessellate"`
	// 	Coordinates string `xml:"coordinates"`
	// } `xml:"LineString"`
}

type KMLParentFolder struct {
	ID string `xml:"id,attr"`
	// Text      string `xml:",chardata"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	// Metadata  KMLMetadata    `xml:"metadata"`
	// Placemark []KMLPlacemark `xml:"Placemark"`
	Folder []KMLFolder `xml:"Folder"`
}

type KMLFolder struct {
	ID string `xml:"id,attr"`
	// Text      string `xml:",chardata"`
	Name      string         `xml:"name"`
	Metadata  KMLMetadata    `xml:"metadata"`
	StyleURL  string         `xml:"styleUrl"`
	Placemark []KMLPlacemark `xml:"Placemark"`
}

type KMLFileSchema struct {
	Name        string                     `xml:"name,attr"`
	ID          string                     `xml:"id,attr"`
	SimpleField []KMLFileSchemaSimpleField `xml:"SimpleField"`
}

type KMLFileSchemaSimpleField struct {
	Type        string `xml:"type,attr"`
	Name        string `xml:"name,attr"`
	DisplayName string `xml:"displayName"`
}

type ExtendedData struct {
	SchemaData SchemaData `xml:"SchemaData"`
	// Data []ExtendedDataData `xml:"Data"`
}

type ExtendedDataData struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

type SchemaData struct {
	SchemaUrl  string       `xml:"schemaUrl,attr"`
	SimpleData []SimpleData `xml:"SimpleData"`
}

type SimpleData struct {
	Value string `xml:",chardata"`
	Name  string `xml:"name,attr"`
}

type Kml struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`
	Lon    float64            `json:"lon" bson:"lon"`
	Lat    float64            `json:"lat" bson:"lat"`
	Type   string             `json:"type" bson:"type"`
	Name   string             `json:"name" bson:"name"`
	// Status int64              `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)
	// Tags   []string           `json:"tags" bson:"tags"`
	// Like   int                `json:"like" bson:"like"`
	// Dlike  int                `json:"dlike" bson:"dlike"`
	//TagsData interface{}        `json:"tagsData" bson:"tags_data"`
	Data   []KmlNodedata `json:"data,omitempty" bson:"data,omitempty"`
	Images []model.Image `json:"images,omitempty" bson:"images,omitempty"`
	// Amenity model.Amenity `json:"amenity,omitempty" bson:"amenity,omitempty"`
	User  model.User `json:"user,omitempty" bson:"user,omitempty"`
	CCode string     `json:"ccode" bson:"ccode"`

	OsmID     string             `json:"osmId" bson:"osm_id"`
	AmenityID primitive.ObjectID `json:"amenityId" bson:"amenity_id"`
	// ReviewsInfo model.ReviewInfo   `json:"reviewinfo" bson:"reviewinfo"`
	// Amenity   []string           `json:"amenity" bson:"amenity"`
	Props     map[string]interface{} `json:"props" bson:"props"`
	CreatedAt time.Time              `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updated_at"`
}

type KmlNodedata struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userId" bson:"user_id"`

	NodeID   primitive.ObjectID `json:"nodeId" bson:"node_id"`
	TagID    primitive.ObjectID `json:"tagId" bson:"tag_id"`
	TagoptID primitive.ObjectID `json:"tagoptId" bson:"tagopt_id"`
	// Name     string                       `json:"name" bson:"name"`
	Tag    model.Tag      `json:"tag" bson:"tag"`
	Tagopt model.Question `json:"tagopt" bson:"tagopt"`
	// Data        model.NodedataData           `json:"data" bson:"data"`
	Title       string                       `json:"title" bson:"title"`
	Description string                       `json:"description" bson:"description"`
	Locale      map[string]map[string]string `json:"locale" bson:"locale"`
	Status      int64                        `json:"status" bson:"status"` // 1 - view, 100 - new, -1 - not view(remove)

	User  model.User            `json:"user,omitempty" bson:"user,omitempty"`
	Audit []model.NodedataAudit `json:"audit,omitempty" bson:"audit,omitempty"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type KMLDescriptionItem struct {
	Key   string      `json:"key" bson:"key"`
	Value interface{} `json:"value" bson:"value"`
}

type KMLParseField struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Lat         string `json:"lat" bson:"lat"`
	Lon         string `json:"lon" bson:"lon"`
	Author      string `json:"author" bson:"author"`
}
