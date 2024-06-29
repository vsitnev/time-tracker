package v1

import (
	"net/http"
	"time-tracker/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "time-tracker/docs"
)

func NewRouter(handler *gin.Engine, services *service.Services) {
	handler.Use(gin.Recovery())

	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "time-tracker",
			"status":  "ok",
		})
	})

	v1 := handler.Group("/api/v1")
	{
		newUserRoutes(v1.Group("/users"), services.User)
		newTaskRoutes(v1.Group("/tasks"), services.Task)
	}
}