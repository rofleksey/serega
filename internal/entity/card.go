package entity

import "time"

const (
	CardStatusTodo  = "todo"
	CardStatusDoing = "doing"
	CardStatusDone  = "done"
)

// Card belongs to the single shared board. Version guards concurrent edits.
type Card struct {
	ID          string
	Title       string
	Description string
	Status      string
	Version     int64
	CreatedBy   User
	UpdatedBy   User
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
