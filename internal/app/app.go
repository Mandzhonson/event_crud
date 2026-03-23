package app

import (
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/logger"
	"calendar/internal/repository"
	"calendar/internal/router"
	"calendar/internal/service"
	"calendar/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() error {
	logger.NewLogger("debug")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	db, err := storage.GetDBConn(ctx, cfg.GetDBString())
	if err != nil {
		return fmt.Errorf("failed to connect db: %w\n", err)
	}
	defer db.Close()
	repo := repository.NewPostgres(db)
	service := service.NewEventService(repo)
	handlers := handlers.NewEventHandler(service)
	router := router.GetRouter(handlers)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTPConfig.Host, cfg.HTTPConfig.Port),
		Handler: router.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Error("failed to start server: ", slog.Any("error", err))
		}
	}()
	slog.Info("server started", slog.String("Host", cfg.HTTPConfig.Host), slog.String("Port", cfg.HTTPConfig.Port))
	<-ctx.Done()
	cancel()
	shtCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shtCtx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server is shutdown succesfully")
	return nil
}
