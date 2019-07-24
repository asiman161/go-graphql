package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"

	"github.com/asiman161/go-graphql/graphql"
	"github.com/asiman161/go-graphql/graphql/dataloaders"
	"github.com/asiman161/go-graphql/localdb"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))


	db := &localdb.LocalDb{}

	rootHandler := dataloaders.DataloaderMiddleware(db, handler.GraphQL(graphql.NewExecutableSchema(
		graphql.NewRootResolvers(db),
	), handler.ComplexityLimit(250)))

	//http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(
	//	graphql.NewRootResolvers(&gorm.DB{}),
	//	), handler.ComplexityLimit(250)))
	http.Handle("/query", rootHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
