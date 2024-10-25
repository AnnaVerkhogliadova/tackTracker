package model

import (
	"github.com/rs/zerolog"
	"time"
)

type Task struct {
	ID          uint64    `db:"task_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      uint64    `db:"status"`
	CreatedAt   time.Time `db:"create_date"`
}

func (t Task) MarshalZerologObject(e *zerolog.Event) {
	e.
		Uint64("task_id", t.ID).
		Str("title", t.Title).
		Str("description", t.Description).
		Uint64("status", t.Status).
		Time("created_at", t.CreatedAt)
}
