package pgdb

import (
	"context"
	"errors"
	"fmt"
	"time"
	"time-tracker/internal/model"
	"time-tracker/internal/repository/repoerr"
	"time-tracker/pkg/postgres"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	maxPaginationLimit     = 10
	defaultPaginationLimit = 10
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepo {
	return &UserRepo{db}
}

type CreateUserInput struct {
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Patronymic     string `json:"patronymic"`
	PassportNumber string `json:"passport_number"`
	Address        string `json:"address"`
}

func (r *UserRepo) CreateUser(ctx context.Context, data CreateUserInput) (int, error) {
	var ID int
	sql, args, _ := r.Builder.Insert("md.users").
		Columns("username", "surname", "patronymic", "passport_number", "address", "created_at").
		Values(data.Name, data.Surname, data.Patronymic, data.PassportNumber, data.Address, time.Now()).
		Suffix("RETURNING id").
		ToSql()

	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("UserRepo.CreateUser - r.Pool.QueryRow: %v", err)
	}
	return ID, nil
}

func (r *UserRepo) GetUser(ctx context.Context, ID int) (model.User, error) {
	var user model.User

	sql, args, err := r.Builder.Select("*").From("md.users").Where("id = ?", ID).ToSql()
	if err != nil {
		return user, fmt.Errorf("UserRepo.GetUser - r.Builder.ToSql: %v", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.PassportNumber, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, repoerr.ErrNotFound
		}
		return user, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) GeUsertByPassportNumber(ctx context.Context, passportNumber string) (model.User, error) {
	var user model.User

	sql, args, err := r.Builder.Select("*").From("md.users").Where("passport_number = ?", passportNumber).ToSql()
	if err != nil {
		return user, fmt.Errorf("UserRepo.GetUser - r.Builder.ToSql: %v", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.PassportNumber, &user.Address, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, repoerr.ErrNotFound
		}
		return user, fmt.Errorf("SourceRepo.GetRouteBySourceID - r.Pool.QueryRow: %v", err)
	}

	return user, nil
}

func (r *UserRepo) ListUsersPagination(ctx context.Context, name, surname, patronymic, passport_number, address string, limit, offset int) ([]model.User, error) {
	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}
	if limit == 0 {
		limit = defaultPaginationLimit
	}

	var whereClauses []squirrel.Sqlizer
	if name != "" {
		whereClauses = append(whereClauses, squirrel.ILike{"username": fmt.Sprintf("%%%s%%", name)})
	}
	if surname != "" {
		whereClauses = append(whereClauses, squirrel.ILike{"surname": fmt.Sprint("%%%s%%", surname)})
	}
	if patronymic != "" {
		whereClauses = append(whereClauses, squirrel.ILike{"patronymic": fmt.Sprint("%%%s%%", patronymic)})
	}
	if passport_number != "" {
		whereClauses = append(whereClauses, squirrel.ILike{"passport_number": fmt.Sprint("%%%s%%", passport_number)})
	}
	if address != "" {
		whereClauses = append(whereClauses, squirrel.ILike{"address": fmt.Sprint("%%%s%%", address)})
	}

	query := r.Builder.Select("*").From("md.users").Limit(uint64(limit)).Offset(uint64(offset))
	if len(whereClauses) > 0 {
		query = query.Where(squirrel.And(whereClauses))
	}
	sql, args, _ := query.ToSql()

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.ListUsersPagination - r.Pool.Query: %v", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		user := model.User{}
		err = rows.Scan(
			&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.PassportNumber, &user.Address, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("UserRepo.ListUsersPagination - rows.Scan: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

type UpdateUserInput struct {
	Name           *string `json:"name"`
	Surname        *string `json:"surname"`
	Patronymic     *string `json:"patronymic"`
	PassportNumber *string `json:"passport_number"`
	Address        *string `json:"address"`
}

func (r *UserRepo) UpdateUser(ctx context.Context, ID int, data UpdateUserInput) error {
	b := r.Builder.Update("md.users").Set("updated_at = ?", time.Now())
	if data.Name != nil {
		b = b.Set("username", *data.Name)
	}
	if data.Surname != nil {
		b = b.Set("surname", *data.Surname)
	}
	if data.Patronymic != nil {
		b = b.Set("patronymic", *data.Patronymic)
	}
	if data.PassportNumber != nil {
		b = b.Set("passport_number", *data.PassportNumber)
	}
	if data.Address != nil {
		b = b.Set("address", *data.Address)
	}
	sql, args, _ := b.Where("id = ?", ID).ToSql()

	_, err := r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo.UpdateUser - r.Pool.Exec: %v", err)
	}

	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int) error {
	sql, args, err := r.Builder.Delete("md.users").Where("id = ?", id).ToSql()
	if err != nil {
		return fmt.Errorf("UserRepo.DeleteUser - r.Builder.ToSql: %v", err)
	}
	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UserRepo.DeleteUser - r.Pool.Exec: %v", err)
	}
	return nil
}
