package v1

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time-tracker/internal/model"
	"time-tracker/internal/repository/pgdb"
	"time-tracker/internal/repository/repoerr"
	"time-tracker/internal/service"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	service service.User
}

func newUserRoutes(handler *gin.RouterGroup, service service.User) {
	r := &UserRoutes{service}
	handler.POST("", r.create)
	handler.GET("", r.getList)
	handler.GET(":id", r.get)
	handler.PATCH(":id", r.update)
	handler.DELETE(":id", r.delete)
}

type createUserInput struct {
	PassportNumber string `json:"passportNumber" binding:"required"`
}
type createUserResponse struct {
	ID int `json:"id"`
}

// @Summary Создание элемента "Пользователь"
// @Description User
// @Tags Users / Пользователи
// @Accept json
// @Produce json
// @Param input body createUserInput true "User input"
// @Success 200 {object} createUserResponse
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/users [post]
func (r *UserRoutes) create(c *gin.Context) {
	var input createUserInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	re, _ := regexp.Compile(`^\d{4} \d{6}$`)
	if !re.MatchString(input.PassportNumber) {
		newErrorResponse(c, http.StatusBadRequest, "invalid passport number")
		return
	}

	id, err := r.service.CreateUser(c, input.PassportNumber)
	if err != nil {
		if errors.Is(err, repoerr.ErrAlreadyExists) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, createUserResponse{
		ID: id,
	})
}

// @Summary Получение элемента "Пользователь"
// @Description Get User
// @Tags Users / Пользователи
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/users/{id} [get]
func (r *UserRoutes) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	item, err := r.service.GetUser(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

type getUserListInput struct {
	Name           string `json:"name,omitempty" form:"name"`
	Surname        string `json:"surname,omitempty" form:"surname"`
	Patronymic     string `json:"patronymic,omitempty" form:"patronymic"`
	PassportNumber string `json:"passportNumber,omitempty" form:"passportNumber"`
	Address        string `json:"address,omitempty" form:"address"`
	Offset         int    `json:"offset,omitempty" form:"offset"`
	Limit          int    `json:"limit,omitempty" form:"limit"`
}

// @Summary Получение списка элементов "Пользователь"
// @Description User list
// @Tags Users / Пользователи
// @Accept json
// @Produce json
// @Param input query getUserListInput true "Filter"
// @Success 200 {array} model.User
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/users [get]
func (r *UserRoutes) getList(c *gin.Context) {
	var input getUserListInput
	if err := c.BindQuery(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if msg, ok := validateUser(input); !ok {
		newErrorResponse(c, http.StatusBadRequest, msg)
		return
	}

	items, err := r.service.ListUsers(c, service.ListUsersFilter{
		Name:           input.Name,
		Surname:        input.Surname,
		Patronymic:     input.Patronymic,
		PassportNumber: input.PassportNumber,
		Address:        input.Address,
		Offset:         input.Offset,
		Limit:          input.Limit,
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if len(items) == 0 {
		items = []model.User{}
	}
	c.JSON(http.StatusOK, items)
}

func validateUser(input getUserListInput) (string, bool) {
	var errs []string
	pattern := regexp.MustCompile(`^[a-zA-Z]+$`)
	if input.Name != "" && (!pattern.MatchString(input.Name) || len(input.Name) > 36) {
		errs = append(errs, "name is invalid")
	}
	if input.Surname != "" && (!pattern.MatchString(input.Surname) || len(input.Surname) > 36) {
		errs = append(errs, "surname is invalid")
	}
	if input.Patronymic != "" && (!pattern.MatchString(input.Patronymic) || len(input.Patronymic) > 36) {
		errs = append(errs, "patronymic is invalid")
	}
	if input.Address != "" && len(input.Address) > 256 {
		errs = append(errs, "address too long")
	}

	re, _ := regexp.Compile(`^\d{4} \d{6}$`)
	if !re.MatchString(input.PassportNumber) {
		errs = append(errs, "passport number is invalid")
	}

	msg := strings.Join(errs, ", ")
	return msg, len(errs) == 0
}

type updateUserInput struct {
	Name           *string `json:"name,omitempty"`
	Surname        *string `json:"surname,omitempty"`
	Patronymic     *string `json:"patronymic,omitempty"`
	PassportNumber *string `json:"passportNumber,omitempty"`
	Address        *string `json:"address,omitempty"`
}
type updateUserResponse struct {
	Success bool `json:"success"`
}

// @Summary Update user
// @Description Update user
// @Tags Users / Пользователи
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body updateUserInput true "User input"
// @Success 200 {object} updateUserResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/users/{id} [patch]
func (r *UserRoutes) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}
	var input updateUserInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if msg, ok := validateUpdateUser(input); !ok {
		newErrorResponse(c, http.StatusBadRequest, msg)
		return
	}
	err = r.service.UpdateUser(c, id, pgdb.UpdateUserInput{
		Name:           input.Name,
		Surname:        input.Surname,
		Patronymic:     input.Patronymic,
		PassportNumber: input.PassportNumber,
		Address:        input.Address,
	})
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, updateUserResponse{
		Success: true,
	})
}

func validateUpdateUser(input updateUserInput) (string, bool) {
	var errs []string
	pattern := regexp.MustCompile(`^[a-zA-Z]+$`)
	if input.Name != nil && (*input.Name != "" && (!pattern.MatchString(*input.Name) || len(*input.Name) > 36)) {
		errs = append(errs, "name is invalid")
	}
	if input.Surname != nil && (*input.Surname != "" && (!pattern.MatchString(*input.Surname) || len(*input.Surname) > 36)) {
		errs = append(errs, "surname is invalid")
	}
	if input.Patronymic != nil && (*input.Patronymic != "" && (!pattern.MatchString(*input.Patronymic) || len(*input.Patronymic) > 36)) {
		errs = append(errs, "patronymic is invalid")
	}
	if input.Address != nil && (*input.Address != "" && len(*input.Address) > 256) {
		errs = append(errs, "address too long")
	}

	re, _ := regexp.Compile(`^\d{4} \d{6}$`)
	if input.PassportNumber != nil && !re.MatchString(*input.PassportNumber) {
		errs = append(errs, "passport number is invalid")
	}

	msg := strings.Join(errs, ", ")
	return msg, len(errs) == 0
}

type deleteUserResponse struct {
	Success bool `json:"success"`
}

// @Summary Delete user
// @Description Delete user
// @Tags Users / Пользователи
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} deleteUserResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/users/{id} [delete]
func (r *UserRoutes) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	err = r.service.DeleteUser(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, deleteUserResponse{
		Success: true,
	})
}
