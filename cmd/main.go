package main

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/http-server/handler/redirect"
	"awesomeProject/internal/http-server/handler/save"
	"awesomeProject/internal/storage"
	"awesomeProject/internal/storage/memory"
	"awesomeProject/internal/storage/postgres"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.Load()

	logger := setupLogger(cfg.Env)
	logger.Info("Server configuration",
		"address", cfg.HTTPServer.Address,
		"timeout", cfg.HTTPServer.Timeout,
		"idle_timeout", cfg.HTTPServer.IdleTimeout,
	)
	logger.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))

	var storage storage.Storage
	var err error

	switch cfg.StorageType {
	case "postgres":
		logger.Info("Storage type", slog.String("type", cfg.StorageType))
		storage, err = postgres.New(cfg.StoragePath)
		if err != nil {
			logger.Error("failt to init storeg", err)
			os.Exit(1)
		}
	case "memory":
		logger.Info("Using memory store")
		storage = memory.New()
	default:
		logger.Error("invalid storage type", slog.String("storage", cfg.StorageType))
		os.Exit(1)
	}

	logger.Info("Storage initialized", slog.String("type", cfg.StorageType))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post(
		"/url",
		save.New(logger, storage),
	)
	router.Get(
		"/{shortUrl}",
		redirect.New(logger, storage),
	)

	logger.Info("starting server", slog.String("address", cfg.Address))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(
				"failed to start server",
				slog.String("address", cfg.Address),
				slog.Any("error", err),
			)
		}
	}()

	logger.Info("server started", slog.String("address", cfg.Address))

	<-done

	logger.Info("stopping server", slog.String("address", cfg.Address))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(
			"failed to stop server",
			slog.String("address", cfg.Address),
			slog.Any("error", err),
		)
		return
	}

	logger.Info("server stopped", slog.String("address", cfg.Address))

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
