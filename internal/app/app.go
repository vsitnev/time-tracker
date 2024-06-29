package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time-tracker/config"
	"time-tracker/internal/repository"
	"time-tracker/internal/service"
	v1 "time-tracker/internal/transport/http/v1"
	"time-tracker/pkg/httpserver"
	"time-tracker/pkg/postgres"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)


// @title Time Tracker
// @version 1.0
// @description A simple time tracker

// @contact.name   Vladislav Sitnev
// @contact.email  sitnevvl@gmail.com

// @host      localhost:8080
// @BasePath  /
func Run() {
	// init Config
	log.Info("Startin application...")
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("error init config", err)
	}

	// init Logger -- logrus
	SetLogger(cfg.Log.Level)

	// init db
	log.Info("Initializing postgres...")
	pg, err := postgres.New(cfg.DSN.Database)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - pgdb.NewServices: %w", err))
	}

	// init Repositories
	log.Info("Initializing repositories...")
	reps := repository.NewRepositories(pg)

	// init Services
	log.Info("Initializing services...")
	deps := service.ServiceDeps{
		Reps: reps,
		ApiURLS: cfg.API,
	}
	services := service.NewServices(deps)

	// Gin handler
	log.Info("Initializing handlers and routes...")
	handler := gin.New()
	v1.NewRouter(handler, services)

	// HTTP server
	log.Info("Starting http server...")
	log.Debugf("Server port: %s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	log.Debugf("Swagger: http://0.0.0.0:%s/swagger/index.html", cfg.HTTP.Port)

	// Waiting signal
	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <- interrupt:
		log.Info("app - Run - signal: " + s.String())
	case err := <- httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
