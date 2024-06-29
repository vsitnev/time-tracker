package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"time-tracker/internal/model"
	"time-tracker/internal/repository"
	"time-tracker/internal/repository/pgdb"
	"time-tracker/internal/repository/repoerr"
)

type UserService struct {
	repo       repository.User
	userApiURl string
}

func NewUserService(repo repository.User, userApiURl string) *UserService {
	return &UserService{
		repo:       repo,
		userApiURl: userApiURl,
	}
}

func (s *UserService) CreateUser(ctx context.Context, passportNumber string) (int, error) {
	_, err := s.repo.GeUsertByPassportNumber(ctx, passportNumber)
	if err == nil {
		return 0, repoerr.ErrAlreadyExists
	}

	passportNums := strings.Split(passportNumber, " ")
	queryParams := fmt.Sprintf("?passportSerie=%s&passportNumber=%s", passportNums[0], passportNums[1])

	client := http.Client{
		Timeout: 4 * time.Second,
	}

	res, err := client.Get(s.userApiURl + "/info" + queryParams)
	if err != nil {
		return 0, fmt.Errorf("UserService.CreateUser - http.Get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("UserService.CreateUser - http.Get: %v", err)
	}

	type userInfoResponse struct {
		Surname    string `json:"surname"`
		Name       string `json:"name"`
		Patronymic string `json:"patronymic"`
		Address    string `json:"address"`
	}
	var info userInfoResponse
	err = json.NewDecoder(res.Body).Decode(&info)
	if err != nil {
		return 0, fmt.Errorf("UserService.CreateUser - pgdb.Unmarshal: %v", err)
	}

	return s.repo.CreateUser(ctx, pgdb.CreateUserInput{
		Name:           info.Name,
		Patronymic:     info.Patronymic,
		PassportNumber: passportNumber,
		Address:        info.Address,
	})
}

func (s *UserService) GetUser(ctx context.Context, ID int) (model.User, error) {
	return s.repo.GetUser(ctx, ID)
}

type ListUsersFilter struct {
	Name           string
	Surname        string
	Patronymic     string
	PassportNumber string
	Address        string
	Limit          int
	Offset         int
}

func (s *UserService) ListUsers(ctx context.Context, filter ListUsersFilter) ([]model.User, error) {
	return s.repo.ListUsersPagination(ctx, filter.Name, filter.Surname, filter.Patronymic, filter.PassportNumber, filter.Address, filter.Limit, filter.Offset)
}

func (s *UserService) UpdateUser(ctx context.Context, ID int, data pgdb.UpdateUserInput) error {
	_, err := s.repo.GetUser(ctx, ID)
	if err != nil {
		return err
	}
	return s.repo.UpdateUser(ctx, ID, data)
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}
