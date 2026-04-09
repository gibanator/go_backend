package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gibanator/go-server/internal/core/logger"
	core_http_server "github.com/gibanator/go-server/internal/core/transport/http/server"
	users_transport_http "github.com/gibanator/go-server/internal/feature/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := logger.NewLogger(logger.LoadConfig())
	if err != nil {
		fmt.Println("failed to init logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Started Application")

	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(nil)
	usersRoutes := usersTransportHTTP.Routes()

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewHTTPServer(
		core_http_server.LoadServerConfig(),
		logger,
	)

	httpServer.RegisterAPIRoutes(apiVersionRouter)

	if err := httpServer.Launch(ctx); err != nil {
		logger.Error("Http server run error:", zap.Error(err))
	}
}
