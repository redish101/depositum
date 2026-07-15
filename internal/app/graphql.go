package app

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/redish101/depositum/pkg/graph"
	"github.com/redish101/depositum/pkg/graph/resolver"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQLServer interface {
	Handlers() http.Handler
	Playground() http.Handler
}

type graphQLServer struct {
	services *Services
}

func (s *graphQLServer) Handlers() http.Handler {
	resolvers := resolver.NewResolver(
		s.services.Library,
		s.services.Shelf,
		s.services.Object,
	)

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: resolvers,
			},
		),
	)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func (s *graphQLServer) Playground() http.Handler {
	return playground.Handler("depositum graphQL playground", "/graphql")
}

func NewGraphQLServer(services *Services) GraphQLServer {
	return &graphQLServer{
		services: services,
	}
}
