package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gptikhomirov/go-rest-prod/internal/core/config"
	core_logger "github.com/gptikhomirov/go-rest-prod/internal/core/logger"
	"github.com/gptikhomirov/go-rest-prod/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/middleware"
	core_http_server "github.com/gptikhomirov/go-rest-prod/internal/core/transport/http/server"
	statistics_repository_postgres "github.com/gptikhomirov/go-rest-prod/internal/features/statistics/repository/postgres"
	statistics_service "github.com/gptikhomirov/go-rest-prod/internal/features/statistics/service"
	statistics_transport_http "github.com/gptikhomirov/go-rest-prod/internal/features/statistics/transport/http"
	tasks_repository_postgres "github.com/gptikhomirov/go-rest-prod/internal/features/tasks/repository/postgres"
	tasks_service "github.com/gptikhomirov/go-rest-prod/internal/features/tasks/service"
	tasks_transport_http "github.com/gptikhomirov/go-rest-prod/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/gptikhomirov/go-rest-prod/internal/features/users/repository/postgres"
	users_service "github.com/gptikhomirov/go-rest-prod/internal/features/users/service"
	users_transport_http "github.com/gptikhomirov/go-rest-prod/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/gptikhomirov/go-rest-prod/docs"
)

// @title Go rest prod API
// @version 1.0
// @description kek
// @host 127.0.0.1:5050
// @BasePath /api/v1
func main() {
	cfg := config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("Failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("app time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool...")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_repository_postgres.NewUsersRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_repository_postgres.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	// for test middleware
	apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
		core_http_server.ApiVersion2,
		core_http_middleware.Test("ROUTER MIDDLEWARE"),
	)
	apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,

		//apiVersionRouterV2,
	)
	httpServer.RegisterSwagger()

	if err = httpServer.Run(ctx); err != nil {
		logger.Error("HTTP Server run error", zap.Error(err))
	}
}
