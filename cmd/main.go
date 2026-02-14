package main

import (
	"calendar/internal/config"
	"calendar/internal/handlers"
	"calendar/internal/repository"
	"calendar/internal/router"
	"calendar/internal/service"
	"calendar/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %s\n", err)
	}
	db, err := storage.GetDBConn(ctx, cfg.GetDBString())
	if err != nil {
		log.Fatalf("failed to connect db: %s\n", err)
	}
	defer db.Close()

	repo := repository.NewPostgres(db)
	service := service.NewEventService(repo)
	handlers := handlers.NewEventHandler(service)
	router := router.GetRouter(handlers)
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.HTTPConfing.Host, cfg.HTTPConfing.Port),
		Handler: router.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %s\n", err)
		}
	}()
	fmt.Println("server started:", fmt.Sprintf("%s:%s", cfg.HTTPConfing.Host, cfg.HTTPConfing.Port))
	<-ctx.Done()
	cancel()
	shtCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shtCtx); err != nil {
		log.Fatalf("failed to shutdown server: %s\n", err)
	}
	fmt.Println("server is shutdown successfully!")
}
