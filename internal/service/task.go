package service

import (
	"context"
	"time"
	"time-tracker/internal/model"
	"time-tracker/internal/repository"
	"time-tracker/internal/repository/pgdb"
)

type TaskService struct {
	repo repository.Task
	userService User
}

func NewTaskService(repo repository.Task, userService User) *TaskService {
	return &TaskService{repo, userService}
}

func (s *TaskService) CreateTask(ctx context.Context, data pgdb.CreateTaskInput) (int, error) {
	return s.repo.CreateTask(ctx, data)
}

func (s *TaskService) GetTask(ctx context.Context, ID int) (model.Task, error) {
	return s.repo.GetTask(ctx, ID)
}

func (s *TaskService) ListTasks(ctx context.Context, userID int, filter pgdb.ListTasksFilter) ([]model.Task, error) {
	_, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListTasks(ctx, userID, filter)
}

func (s *TaskService) CompleteTask(ctx context.Context, ID int) error {
	task, err := s.repo.GetTask(ctx, ID)
	if err != nil {
		return err
	}

	if task.Completed {
		return ErrTaskAlreadyCompleted
	}

	return s.repo.UpdateTask(ctx, ID, pgdb.UpdateTaskInput{
		Completed: true,
		UpdatedAt: time.Now(),
		Duration:  int(time.Since(task.CreatedAt).Minutes()),
	})
}