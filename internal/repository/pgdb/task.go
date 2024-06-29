package pgdb

import (
	"context"
	"errors"
	"fmt"
	"time"
	"time-tracker/internal/model"
	"time-tracker/internal/repository/repoerr"
	"time-tracker/pkg/postgres"

	"github.com/jackc/pgx/v5"
)

type TaskRepo struct {
	*postgres.Postgres
}

func NewTaskRepo(db *postgres.Postgres) *TaskRepo {
	return &TaskRepo{db}
}

type CreateTaskInput struct {
	UserID      int
	Description string
}

func (r *TaskRepo) CreateTask(ctx context.Context, data CreateTaskInput) (int, error) {
	var ID int
	sql, args, _ := r.Builder.Insert("md.tasks").
		Columns("user_id", "description", "completed", "duration", "created_at").
		Values(data.UserID, data.Description, false, 0, time.Now()).
		Suffix("RETURNING id").
		ToSql()

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("TaskRepo.CreateTask - r.Pool.QueryRow: %v", err)
	}
	return ID, nil
}

func (r *TaskRepo) GetTask(ctx context.Context, ID int) (model.Task, error) {
	var task model.Task

	sql, args, _ := r.Builder.Select("*").From("md.tasks").Where("id = ?", ID).ToSql()

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&task.ID, &task.UserID, &task.Description, &task.Duration, &task.Completed,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return task, repoerr.ErrNotFound
		}
		return task, fmt.Errorf("TaskRepo.GetTask - r.Pool.QueryRow: %v", err)
	}
	return task, nil
}

type ListTasksFilter struct {
	DateTo   time.Time
	DateFrom time.Time
}

func (r *TaskRepo) ListTasks(ctx context.Context, userID int, filter ListTasksFilter) ([]model.Task, error) {
	sql, args, _ := r.Builder.Select("*").From("md.tasks").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, filter.DateFrom, filter.DateTo).
		OrderBy("duration DESC").
		ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("TaskRepo.ListTasks - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.Description, &task.Duration, &task.Completed,
			&task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("TaskRepo.ListTasks - rows.Scan: %v", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

type UpdateTaskInput struct {
	Completed bool
	Duration  int
	UpdatedAt time.Time
}

func (r *TaskRepo) UpdateTask(ctx context.Context, ID int, data UpdateTaskInput) error {
	sql, args, _ := r.Builder.Update("md.tasks").
		SetMap(map[string]interface{}{
			"completed":  data.Completed,
			"duration":   data.Duration,
			"updated_at": data.UpdatedAt,
		}).
		Where("id = ?", ID).
		ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("TaskRepo.UpdateTask - r.Pool.Exec: %v", err)
	}
	return nil
}
