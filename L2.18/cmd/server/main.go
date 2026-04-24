package main

import (
	"calendar/internal/config"
	internalhttp "calendar/internal/transport/http"
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"calendar/internal/repository"
	"calendar/internal/service"
)

func main() {
	cfg := config.MustLoad()

	logger, err := setupLogger(cfg.Env)
	if err != nil {
		log.Fatalf("logger setup: %s", err)
	}
	logger.Info("Logger started", slog.String("env", cfg.Env))
	logger.Debug("Debug messages enabled")

	repo := repository.NewRepository()

	svc := service.NewCalendarService(repo)

	handler := internalhttp.NewHandler(svc, logger)

	router := internalhttp.NewRouter(handler)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server...", slog.String("address", cfg.Address))

		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	sig := <-stopCh

	logger.Info("Stopping server...", slog.String("signal", sig.String()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", slog.Any("error", err))
	} else {
		logger.Info("Server stopped gracefully")
	}
}

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger
	switch env {
	case config.EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		return nil, fmt.Errorf("invalid env: %s", env)
	}

	return logger, nil
}
