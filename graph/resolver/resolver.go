package resolver

import (
	"github.com/mikalai2006/swapland-api/graph/loaders"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB         *mongo.Database
	Repo       *repository.Repositories
	TagsLoader *loaders.Loaders
}
