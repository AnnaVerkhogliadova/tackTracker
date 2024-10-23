package controller

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"taskTracker/driver"
	"taskTracker/model"
)

type taskController struct {
	taskDriver driver.ITasks
}

func NewController(
	taskDriver driver.ITasks,
) (ITaskController, error) {
	return &taskController{
		taskDriver: taskDriver,
	}, nil
}

func (t taskController) Create(ctx context.Context, task *model.Task) (uint64, error) {
	lgr := zerolog.Ctx(ctx).With().
		Dict("request", zerolog.Dict().
			Interface("task", task)).
		Logger()

	taskId, err := t.taskDriver.Create(ctx, task)
	if err != nil {
		lgr.Error().
			Err(err).
			Msg("error creating Task")

		return 0, fmt.Errorf("error creating task: %w", err)
	}

	return taskId, nil
}

func (t taskController) SetStatus(ctx context.Context, taskId uint64, status *uint64) error {
	err := t.taskDriver.SetStatus(
		ctx,
		taskId,
		status)

	if err != nil {
		return err
	}

	return nil
}

func (t taskController) Get(ctx context.Context, taskId uint64) (*model.Task, error) {
	task, err := t.taskDriver.Get(
		ctx,
		taskId,
	)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t taskController) Delete(ctx context.Context, taskId uint64) error {
	err := t.taskDriver.Delete(
		ctx,
		taskId)

	if err != nil {
		return err
	}

	return nil
}

func (t taskController) GetList(ctx context.Context, status *uint64) ([]*model.Task, error) {
	tasks, err := t.taskDriver.GetList(
		ctx,
		status)

	if err != nil {
		return nil, err
	}

	return tasks, nil
}
