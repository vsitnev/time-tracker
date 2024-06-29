package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"time-tracker/internal/repository/pgdb"
	"time-tracker/internal/repository/repoerr"
	"time-tracker/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskRoutes struct {
	service service.Task
}

func newTaskRoutes(handler *gin.RouterGroup, service service.Task) {
	r := &TaskRoutes{service}
	handler.POST("", r.create)
	handler.GET("", r.getList)
	handler.GET(":id", r.get)
	handler.POST(":id/complete", r.complete)
}

type getTaskListInput struct {
	UserID   int       `json:"userId" form:"userId" binding:"required"`
	DateTo   time.Time `json:"dateTo" time_format:"2006-01-02T15:04:05Z07:00" form:"dateTo" binding:"required"`
	DateFrom time.Time `json:"dateFrom" time_format:"2006-01-02T15:04:05Z07:00" form:"dateFrom" binding:"required"`
}

// @Summary Получение списка элементов "Задача"
// @Description Task list
// @Tags Tasks / Задачи
// @Accept json
// @Produce json
// @Param input query getTaskListInput true "Filter"
// @Success 200 {array} model.Task
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/tasks [get]
func (r *TaskRoutes) getList(c *gin.Context) {
	var input getTaskListInput
	if err := c.BindQuery(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if condition := input.DateFrom.After(input.DateTo); condition {
		newErrorResponse(c, http.StatusBadRequest, "invalid date range")
		return
	}

	items, err := r.service.ListTasks(c, input.UserID, pgdb.ListTasksFilter{
		DateFrom: input.DateFrom,
		DateTo:   input.DateTo,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

// @Summary Получение элемента "Задача"
// @Description Get Task
// @Tags Tasks / Задачи
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} model.Task
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/tasks/{id} [get]
func (r *TaskRoutes) get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}
	item, err := r.service.GetTask(c, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

type createTaskInput struct {
	UserID      int    `json:"userId" form:"userId" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
}

type createTaskResponse struct {
	ID int `json:"id"`
}

// @Summary Создание элемента "Задача"
// @Description Create Task
// @Tags Tasks / Задачи
// @Accept json
// @Produce json
// @Param input body createTaskInput true "Task input"
// @Success 200 {object} createTaskResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/tasks [post]
func (r *TaskRoutes) create(c *gin.Context) {
	var input createTaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := r.service.CreateTask(c, pgdb.CreateTaskInput{
		UserID:      input.UserID,
		Description: input.Description,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, createTaskResponse{
		ID: id,
	})
}

type completeTaskResponse struct {
	Success bool `json:"success"`
}

// @Summary Завершение задачи
// @Description Complete specified task
// @Tags Tasks / Задачи
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} completeTaskResponse
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/v1/tasks/{id}/complete [post]
func (r *TaskRoutes) complete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}
	err = r.service.CompleteTask(c, id)
	if err != nil {
		if errors.Is(err, repoerr.ErrNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, service.ErrTaskAlreadyCompleted) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, completeTaskResponse{
		Success: true,
	})
}
