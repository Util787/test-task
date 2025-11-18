package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Util787/test-task/internal/adapters/rest"
	"github.com/Util787/test-task/internal/adapters/storage"
	"github.com/Util787/test-task/internal/config"
	"github.com/Util787/test-task/internal/logger/slogpretty"
	"github.com/Util787/test-task/internal/usecase"
)

// для упрощения везде сделал context.Background()
func main() {
	cfg := config.MustLoadConfig()

	log := slogpretty.NewPrettyLogger(os.Stdout, slog.LevelDebug)

	postgreStorage := storage.MustInitPostgres(context.Background(), cfg.PostgresConfig)
	log.Debug("Postgres initialized", slog.String("host", cfg.PostgresConfig.Host), slog.Int("port", cfg.PostgresConfig.Port), slog.String("db_name", cfg.PostgresConfig.DbName))

	sortUsecase := usecase.NewSortUsecase(&postgreStorage)

	serv := rest.NewRestServer(log, cfg.HTTPServerConfig, sortUsecase)

	go func() {
		log.Info("HTTP server start", slog.String("host", cfg.HTTPServerConfig.Host), slog.Int("port", cfg.HTTPServerConfig.Port))
		if err := serv.Run(); err != nil {
			log.Error("HTTP server error", slog.String("error", err.Error()))
		}
	}()

	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	log.Info("Shutting down gracefully...")

	log.Info("Shutting down server")
	if err := serv.Shutdown(context.Background()); err != nil {
		log.Error("HTTP server shutdown error", slog.String("error", err.Error()))
	}

	log.Info("Shutting down postgres")
	postgreStorage.Shutdown()

	log.Info("Shutdown complete")

}
