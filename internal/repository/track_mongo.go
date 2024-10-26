package repository

import (
	"context"
	"time"

	"github.com/mikalai2006/swapland-api/internal/config"
	"github.com/mikalai2006/swapland-api/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrackMongo struct {
	db   *mongo.Database
	i18n config.I18nConfig
}

func NewTrackMongo(db *mongo.Database, i18n config.I18nConfig) *TrackMongo {
	return &TrackMongo{db: db, i18n: i18n}
}

func (r *TrackMongo) FindTrack(params domain.RequestParams) (domain.Response[domain.Track], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Track
	var response domain.Response[domain.Track]
	filter, opts, err := CreateFilterAndOptions(params)
	if err != nil {
		return domain.Response[domain.Track]{}, err
	}

	cursor, err := r.db.Collection(tblTrack).Find(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Track, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblTrack).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Track]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TrackMongo) GetAllTrack(params domain.RequestParams) (domain.Response[domain.Track], error) {
	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	var results []domain.Track
	var response domain.Response[domain.Track]
	pipe, err := CreatePipeline(params, &r.i18n)
	if err != nil {
		return domain.Response[domain.Track]{}, err
	}

	cursor, err := r.db.Collection(tblTrack).Aggregate(ctx, pipe) // Find(ctx, params.Filter, opts)
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if er := cursor.All(ctx, &results); er != nil {
		return response, er
	}

	resultSlice := make([]domain.Track, len(results))
	// for i, d := range results {
	// 	resultSlice[i] = d
	// }
	copy(resultSlice, results)

	count, err := r.db.Collection(tblTrack).CountDocuments(ctx, bson.M{})
	if err != nil {
		return response, err
	}

	response = domain.Response[domain.Track]{
		Total: int(count),
		Skip:  int(params.Options.Skip),
		Limit: int(params.Options.Limit),
		Data:  resultSlice,
	}
	return response, nil
}

func (r *TrackMongo) CreateTrack(userID string, track *domain.Track) (*domain.Track, error) {
	var result *domain.Track

	collection := r.db.Collection(tblTrack)

	ctx, cancel := context.WithTimeout(context.Background(), MongoQueryTimeout)
	defer cancel()

	userIDPrimitive, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	if !track.CreatedAt.IsZero() {
		createdAt = track.CreatedAt
	}

	newTrack := domain.Track{
		Lon:       track.Lon,
		Lat:       track.Lat,
		UserID:    userIDPrimitive,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
	}

	res, err := collection.InsertOne(ctx, newTrack)
	if err != nil {
		return nil, err
	}

	err = r.db.Collection(tblTrack).FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
