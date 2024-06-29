package service

import (
	"context"
	"time-tracker/config"
	"time-tracker/internal/model"
	"time-tracker/internal/repository"
	"time-tracker/internal/repository/pgdb"
)

type User interface {
	CreateUser(ctx context.Context, passportNumber string) (int, error)
	GetUser(ctx context.Context, ID int) (model.User, error)
	ListUsers(ctx context.Context, filter ListUsersFilter) ([]model.User, error)
	UpdateUser(ctx context.Context, ID int, data pgdb.UpdateUserInput) error
	DeleteUser(ctx context.Context, id int) error
}

type Task interface {
	CreateTask(ctx context.Context, data pgdb.CreateTaskInput) (int, error)
	GetTask(ctx context.Context, ID int) (model.Task, error)
	ListTasks(ctx context.Context, userID int, filter pgdb.ListTasksFilter) ([]model.Task, error)
	CompleteTask(ctx context.Context, ID int) error
}

type Services struct {
	User
	Task
}

type ServiceDeps struct {
	Reps    *repository.Repositories
	ApiURLS config.API
}

func NewServices(deps ServiceDeps) *Services {
	userService := NewUserService(deps.Reps, deps.ApiURLS.UserApiURl)
	return &Services{
		User: userService,
		Task: NewTaskService(deps.Reps, userService),
	}
}
