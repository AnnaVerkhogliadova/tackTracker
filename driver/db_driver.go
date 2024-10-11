package driver

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	//"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
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

func (d *dbDriver) Create(ctx context.Context, task *model.Task) (uint64, error) {
	tx, err := d.rwdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("error creating task: %w", err)
	}

	var taskId uint64

	err = tx.QueryRow(ctx, queryCreateTask,
		task.Title, task.Description, task.Status,
	).Scan(&taskId)

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("error executing query")
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error creating task in db: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return taskId, nil
}

func (d *dbDriver) Get(ctx context.Context, taskId uint64) (*model.Task, error) {
	tx, err := d.rwdb.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}

	rows, err := tx.Query(
		ctx,
		queryGet,
		taskId)

	if err != nil {
		return nil, fmt.Errorf("error Get in db: %w", err)
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Task])
	if err != nil {
		return nil, fmt.Errorf("errorCollectRows for Get: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("not found: %w", err)
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
		return fmt.Errorf("error deleting Sip: %w", err)
	}

	return nil
}
