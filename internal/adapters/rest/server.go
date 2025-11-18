package rest

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Util787/test-task/internal/config"
)

const ( // чтобы не загромождать конфиг
	defaultReadTimeout       = 5 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
)

type Server struct {
	httpServer *http.Server
}

func NewRestServer(log *slog.Logger, config config.HTTPServerConfig, sortUsecase sortUsecase) Server {
	handler := Handler{
		log:         log,
		sortUsecase: sortUsecase,
	}

	httpServer := &http.Server{
		Addr:              config.Host + ":" + strconv.Itoa(config.Port),
		Handler:           handler.InitRoutes(),
		MaxHeaderBytes:    1 << 20, // 1 MB
		ReadHeaderTimeout: defaultReadHeaderTimeout,
		WriteTimeout:      defaultWriteTimeout,
		ReadTimeout:       defaultReadTimeout,
	}

	return Server{
		httpServer: httpServer,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
