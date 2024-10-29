package driver

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"taskTracker/errors"
	"taskTracker/model"
)

type dbDriver struct {
	rwdb *pgxpool.Pool
	rdb  *pgxpool.Pool
	qb   *squirrel.StatementBuilderType
}

func NewDbDriver(rwdb *pgxpool.Pool, rdb *pgxpool.Pool) (ITasks, error) {
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &dbDriver{
		rwdb: rwdb,
		rdb:  rdb,
		qb:   &qb,
	}, nil
}

func (d *dbDriver) Create(ctx context.Context, task *model.Task) (*model.Task, error) {
	tx, err := d.rwdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}

	var taskId uint64

	err = tx.QueryRow(ctx, queryCreateTask,
		task.Title, task.Description, task.Status,
	).Scan(&taskId)

	task = &model.Task{
		ID:          taskId,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error executing query")
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error creating task in db: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return task, nil
}

func (d *dbDriver) SetStatus(ctx context.Context, taskId uint64, status *uint64) error {
	row, err := d.rwdb.Query(
		ctx,
		querySetStatus,
		taskId, status)

	defer row.Close()
	if err != nil {
		return fmt.Errorf("error set status: %w", err)
	}

	return nil
}

func (d *dbDriver) Get(ctx context.Context, taskId uint64) (*model.Task, error) {
	tx, err := d.rwdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}

	row, err := tx.Query(
		ctx,
		queryGet,
		taskId)

	if err != nil {
		return nil, fmt.Errorf("error Get in db: %w", err)
	}

	results, err := pgx.CollectRows(row, pgx.RowToStructByName[model.Task])
	if err != nil {
		return nil, fmt.Errorf("errorCollectRows for Get: %w", err)
	}

	if len(results) == 0 {
		return nil, errors.NewErrTaskNotFound(taskId)
	}

	return &results[0], nil
}

func (d *dbDriver) Delete(ctx context.Context, taskId uint64) error {
	row, err := d.rwdb.Query(
		ctx,
		queryDelete,
		taskId)

	defer row.Close()
	if err != nil {
		return fmt.Errorf("error deleting Task: %w", err)
	}

	return nil
}

func (d *dbDriver) GetList(ctx context.Context, status *uint64) ([]*model.Task, error) {
	row, err := d.rwdb.Query(
		ctx,
		queryGetList,
		status)

	if err != nil {
		return nil, fmt.Errorf("error Get in db: %w", err)
	}

	results, err := pgx.CollectRows(row, pgx.RowToStructByName[model.Task])
	if err != nil {
		return nil, fmt.Errorf("errorCollectRows for GetList: %w", err)
	}

	tasks := make([]*model.Task, len(results))
	for i := range results {
		task := results[i]
		tasks[i] = &task
	}

	return tasks, nil
}

func (d *dbDriver) CreateSubTask(ctx context.Context, subTask *model.SubTask) (*model.SubTask, error) {
	tx, err := d.rwdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}

	var exists bool
	err = tx.QueryRow(ctx, queryExistTaskId, subTask.TaskID).Scan(&exists)
	if err != nil {
		return nil, nil
	}

	if !exists {
		tx.Rollback(ctx)
		return nil, errors.NewErrTaskNotFound(subTask.TaskID)
	}

	var subTaskId uint64
	err = tx.QueryRow(ctx, queryCreateSubTask,
		subTask.TaskID,
		subTask.Title,
		subTask.Description,
		subTask.Status,
	).Scan(&subTaskId)

	subTask = &model.SubTask{
		ID:          subTaskId,
		TaskID:      subTask.TaskID,
		Title:       subTask.Title,
		Description: subTask.Description,
		Status:      subTask.Status,
		CreatedAt:   subTask.CreatedAt}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error executing query")
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error creating task in db: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return subTask, nil
}
