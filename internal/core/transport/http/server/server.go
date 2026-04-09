package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gibanator/go-server/internal/core/logger"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	cfg Config
	log *logger.Logger
}

func NewHTTPServer(cfg Config, log *logger.Logger) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		cfg: cfg,
		log: log,
	}
}

func (s *HTTPServer) RegisterAPIRoutes(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (s *HTTPServer) Launch(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("HTTP Server started:", zap.String("addr", s.cfg.Addr))
		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("Http error: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("Server shutdown. . .")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("server shutdown with error: %w", err)
		}

		s.log.Warn("Server Stopped gracefully")
	}
	return nil
}
