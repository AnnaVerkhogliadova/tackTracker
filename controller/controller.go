package controller

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"taskTracker/driver"
	"taskTracker/model"
)

type taskController struct {
	taskDriver driver.ITasks
}

func NewController(taskDriver driver.ITasks) (ITaskController, error) {
	return &taskController{
		taskDriver: taskDriver,
	}, nil
}

func (t taskController) Create(ctx context.Context, task *model.Task) (uint64, error) {
	createdTask, err := t.taskDriver.Create(ctx, task)
	if err != nil {
		return 0, fmt.Errorf("error creating task: %w", err)
	}

	logger := log.With().Object("task", createdTask).Logger()
	logger.Info().Msg("Create result")
	return createdTask.ID, nil
}

func (t taskController) SetStatus(ctx context.Context, taskId uint64, status *uint64) error {
	err := t.taskDriver.SetStatus(
		ctx,
		taskId,
		status)

	if err != nil {
		return err
	}

	logger := log.With().Int("status", int(*status)).Logger()
	logger.Info().Msg("SetStatus result")
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

	log.Info().Msg("Get result")
	return task, nil
}

func (t taskController) Delete(ctx context.Context, taskId uint64) error {
	err := t.taskDriver.Delete(
		ctx,
		taskId)

	if err != nil {
		return err
	}

	logger := log.With().Int("task_id", int(taskId)).Logger()
	logger.Info().Msg("Delete result")
	return nil
}

func (t taskController) GetList(ctx context.Context, status *uint64) ([]*model.Task, error) {
	tasks, err := t.taskDriver.GetList(
		ctx,
		status)

	if err != nil {
		return nil, err
	}

	log.Info().Msg("GetList result")
	return tasks, nil
}
