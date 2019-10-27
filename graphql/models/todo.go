package models

import "time"

type Todo struct {
	ID     string
	Text   string
	Done   bool
	UserID string
	Time   time.Time
}

func (Todo) IsUserTodo() {}
