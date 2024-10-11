package model

import "time"

type Task struct {
	ID          uint64    `db:"task_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Status      uint64    `db:"status"`
	CreatedAt   time.Time `db:"create_date"`
}
