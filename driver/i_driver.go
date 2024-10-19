package driver

import (
	"context"
	"taskTracker/model"
)

type ITasks interface {
	Create(ctx context.Context, task *model.Task) (uint64, error)
	Get(ctx context.Context, taskId uint64) (*model.Task, error)
	Delete(ctx context.Context, taskId uint64) error
	GetList(ctx context.Context, status *uint64) ([]*model.Task, error)
}
