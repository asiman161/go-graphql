package dataloaders

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/asiman161/go-graphql/graphql/models"
	"github.com/asiman161/go-graphql/localdb"
)

const LoaderKey = "LoaderCtx"

type Loads struct {
	UserLoader UserLoader
	TodoLoader TodoLoader
}

func DataloaderMiddleware(db *localdb.LocalDb, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLoader := UserLoader{
			maxBatch: 100,
			wait:     1 * time.Millisecond,
			fetch: func(ids []string) ([]*models.User, []error) {
				if len(ids) > 1 {
					fmt.Printf("get many users: %d\n", len(ids))
				} else {
					fmt.Println("get one user")
				}

				var res []*models.User

				for i := range ids {
					for j := range db.Users {
						if db.Users[j].ID == ids[i] {
							res = append(res, db.Users[i])
							break
						}
					}
				}

				if len(res) == 0 {
					return nil, []error{errors.New("users not found")}
				}
				return res, nil
			},
		}

		todoLoader := TodoLoader{
			maxBatch: 100,
			wait:     1 * time.Millisecond,
			fetch: func(ids []string) ([][]*models.Todo, []error) {
				if len(ids) > 1 {
					fmt.Printf("get many todos: %d\n", len(ids))
				} else {
					fmt.Println("get one todo")
				}

				var res [][]*models.Todo

				for i := range ids {
					var tds []*models.Todo
					for j := range db.Todos {
						if db.Todos[j].UserID == ids[i] {
							tds = append(tds, db.Todos[i])
						}
					}
					res = append(res, tds)
				}

				if len(res) == 0 {
					return nil, []error{errors.New("todos not found")}
				}

				return res, nil
			},
		}
		ctx := context.WithValue(r.Context(), LoaderKey, &Loads{UserLoader: userLoader, TodoLoader:todoLoader})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
