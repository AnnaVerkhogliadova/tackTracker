package driver

import (
	"context"
	"taskTracker/model"
)

type ITasks interface {
	Create(ctx context.Context, task *model.Task) (*model.Task, error)
	SetStatus(ctx context.Context, taskId uint64, status *uint64) error
	SetSubTaskStatus(ctx context.Context, subTaskId uint64, status *uint64) error
	Get(ctx context.Context, taskId uint64) (*model.Task, error)
	Delete(ctx context.Context, taskId uint64) error
	DeleteSubTask(ctx context.Context, subTaskId uint64) error
	GetList(ctx context.Context, status *uint64) ([]*model.Task, error)
	CreateSubTask(ctx context.Context, subTask *model.SubTask) (*model.SubTask, error)
}
