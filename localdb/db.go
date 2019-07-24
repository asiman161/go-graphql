package localdb

import "github.com/asiman161/go-graphql/graphql/models"

type LocalDb struct {
	Users []*models.User
	Todos []*models.Todo
}
