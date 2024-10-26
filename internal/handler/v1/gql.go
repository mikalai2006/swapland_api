package v1

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/mikalai2006/swapland-api/graph/generated"
	"github.com/mikalai2006/swapland-api/graph/loaders"
	"github.com/mikalai2006/swapland-api/graph/resolver"
	"github.com/mikalai2006/swapland-api/internal/middleware"
	"github.com/mikalai2006/swapland-api/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

// Defining the Graphql handler
func graphqlHandler(mongoDB *mongo.Database, repositories *repository.Repositories) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	tagsLoader := loaders.NewLoaders(mongoDB, repositories)
	c := generated.Config{Resolvers: &resolver.Resolver{
		DB:         mongoDB,
		Repo:       repositories,
		TagsLoader: tagsLoader,
	}}
	c.Directives.Auth = middleware.GetAuth

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.SetQueryCache(graphql.NoCache{})
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/api/v1/gql/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *HandlerV1) registerGql(router *gin.RouterGroup) {
	router.Use(middleware.GinContextToContextMiddleware())
	var gql = router.Group("/gql")
	gql.Use(loaders.Middleware(h.db, h.repositories))
	gql.GET("/", playgroundHandler())
	gql.POST("/query", middleware.SetUserIdentityGraphql, graphqlHandler(h.db, h.repositories))

}
