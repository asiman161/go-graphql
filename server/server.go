package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	graphql2 "github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"

	"github.com/asiman161/go-graphql/graphql"
	"github.com/asiman161/go-graphql/graphql/dataloaders"
	"github.com/asiman161/go-graphql/graphql/models"
	"github.com/asiman161/go-graphql/localdb"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("auth")
			role := r.Header.Get("auth")
			fmt.Printf("MIDDLEWARE[auth]: {auth: \"%s\", role: \"%s\"} \n", auth, role)
			ctx := r.Context()
			ctx = context.WithValue(ctx, "auth", auth)
			ctx = context.WithValue(ctx, "role", role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	db := &localdb.LocalDb{}


	r.Handle("/", handler.Playground("GraphQL playground", "/query"))
	config := graphql.NewRootResolvers(db)
	config.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql2.Resolver, role models.Role) (res interface{}, err error) {
		fmt.Printf("DIRECTIVE[role]: {auth: \"%s\", role: \"%s\"} \n", ctx.Value("auth"), ctx.Value("role"))
		return next(ctx)
	}

	config.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql2.Resolver) (res interface{}, err error) {
		fmt.Printf("DIRECTIVE[auth]: {auth: \"%s\", role: \"%s\"} \n", ctx.Value("auth"), ctx.Value("role"))
		return next(ctx)
	}

	rootHandler := dataloaders.DataloaderMiddleware(db, handler.GraphQL(graphql.NewExecutableSchema(
		config,
	), handler.ComplexityLimit(250)))

	r.Handle("/query", rootHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
