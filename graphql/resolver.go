package graphql

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/asiman161/go-graphql/graphql/dataloaders"
	"github.com/asiman161/go-graphql/graphql/models"
	"github.com/asiman161/go-graphql/localdb"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	db *localdb.LocalDb
}

func NewRootResolvers(db *localdb.LocalDb) Config {
	db.Users = []*models.User{
		{ID: "1", Name: "Alex", Email: "alex@google.com"},
		{ID: "2", Name: "John", Email: "john@google.com"},
	}
	db.Todos = []*models.Todo{
		{ID: "1", Text: "Alex first message", Done: true, UserID: "1", Time: time.Now()},
		{ID: "2", Text: "Alex second message", Done: false, UserID: "1", Time: time.Now()},

		{ID: "3", Text: "John random msg", Done: false, UserID: "2", Time: time.Now()},
	}
	c := Config{
		Resolvers: &Resolver{
			db: db,
		},
	}

	// Complexity
	countComplexity := func(childComplexity int, limit *int, offset *int) int {
		return *limit * childComplexity
	}
	c.Complexity.Query.Todos = countComplexity
	c.Complexity.Query.Users = countComplexity

	// Schema Directive
	//c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	//	ctxUserID := ctx.Value(UserIDCtxKey)
	//	if ctxUserID != nil {
	//		return next(ctx)
	//	} else {
	//		return nil, errors.UnauthorisedError
	//	}
	//}
	return c
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input models.NewUser) (*models.User, error) {
	user := &models.User{
		Name:  input.Name,
		ID:    fmt.Sprintf("T%d", rand.Intn(1000)),
		Email: input.Email,
	}
	r.db.Users = append(r.db.Users, user)
	return user, nil
}

func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*models.User, error) {
	return r.db.Users, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	for i := range r.db.Users {
		if r.db.Users[i].ID == id {
			return r.db.Users[i], nil
		}
	}

	return nil, errors.New("user not found")
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input models.NewTodo) (*models.Todo, error) {
	for i := range r.db.Users {
		if r.db.Users[i].ID == input.UserID {
			todo := &models.Todo{
				Text:   input.Text,
				ID:     fmt.Sprintf("T%d", rand.Intn(1000)),
				UserID: input.UserID,
				Time:   time.Now(),
			}
			r.db.Todos = append(r.db.Todos, todo)
			return todo, nil
		}
	}
	return &models.Todo{}, errors.New("user not found")
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input models.UpdateTodo) (*models.Todo, error) {
	for i := range r.db.Todos {
		if r.db.Todos[i].ID == input.TodoID {
			r.db.Todos[i].Done = input.Done
			return r.db.Todos[i], nil
		}
	}
	return nil, errors.New("not Found")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context, limit *int, offset *int) ([]*models.Todo, error) {
	return r.db.Todos, nil
}

func (r *queryResolver) Todo(ctx context.Context, id string) (*models.Todo, error) {
	for i := range r.db.Todos {
		if r.db.Todos[i].ID == id {
			return r.db.Todos[i], nil
		}
	}
	return nil, errors.New("not Found")
}

type userResolver struct{ *Resolver }

func (r *userResolver) Todos(ctx context.Context, user *models.User, limit *int, offset *int) ([]*models.Todo, error) {
	return ctx.Value(dataloaders.LoaderKey).(*dataloaders.Loads).TodoLoader.Load(user.ID)
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, todo *models.Todo, limit *int, offset *int) (*models.User, error) {
	return ctx.Value(dataloaders.LoaderKey).(*dataloaders.Loads).UserLoader.Load(todo.UserID)
}
