package repository

import (
	"context"
	"time-tracker/internal/model"
	"time-tracker/internal/repository/pgdb"
	"time-tracker/pkg/postgres"
)

type User interface{
	CreateUser(ctx context.Context, data pgdb.CreateUserInput) (int, error)
	GetUser(ctx context.Context, ID int) (model.User, error)
	GeUsertByPassportNumber(ctx context.Context, passportNumber string) (model.User, error)
	ListUsersPagination(ctx context.Context, name, surname, patronymic, passport_number, address string, limit, offset int) ([]model.User, error)
	UpdateUser(ctx context.Context, ID int, data pgdb.UpdateUserInput) error
	DeleteUser(ctx context.Context, id int) error
}

type Task interface{
	CreateTask(ctx context.Context, data pgdb.CreateTaskInput) (int, error)
	GetTask(ctx context.Context, ID int) (model.Task, error)
	ListTasks(ctx context.Context, userID int, filter pgdb.ListTasksFilter) ([]model.Task, error)
	UpdateTask(ctx context.Context, ID int, data pgdb.UpdateTaskInput) error
}

type Repositories struct {
	User
	Task
}

func NewRepositories(db *postgres.Postgres) *Repositories {
	return &Repositories{
		User: pgdb.NewUserRepo(db),
		Task: pgdb.NewTaskRepo(db),
	}
}