package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"scm.x5.ru/x5m/go-backend/template/internal/api"
	"scm.x5.ru/x5m/go-backend/template/internal/config"
	"scm.x5.ru/x5m/go-backend/template/internal/middlewares"
	"scm.x5.ru/x5m/go-backend/template/internal/services"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"scm.x5.ru/x5m/go-backend/packages/httputils"
	"scm.x5.ru/x5m/go-backend/packages/metrics"
	db "scm.x5.ru/x5m/go-backend/packages/postgres"
	"scm.x5.ru/x5m/go-backend/packages/zlogger"
)

var defaultLogger = zlogger.Logger

func main() {
	// Logging setup
	defaultLogger.Info().Msg("start template service")

	logger := defaultLogger
	ctx := logger.WithContext(context.Background())

	// Configuration load & setup
	panicOnErr(config.Initialize(ctx))
	cfg := config.Get()

	// Swagger setup
	log.Info().Msg("start swagger")
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Error().Err(err).Msg("Error loading swagger spec")
		return
	}

	swagger.Servers = nil

	// Storage & Communication setup
	_ = initializeDB(ctx, &cfg)

	// Server setup
	pingSrv := services.NewPingService(&cfg)

	deps := api.ControllerDeps{
		PingSrv: pingSrv,
	}
	controller := api.NewController(&cfg, deps)

	runner, _ := errgroup.WithContext(ctx)

	runner.Go(NewHTTPServer(ctx, &cfg, controller, swagger))

	runner.Go(func() error {
		httpServer, err := metrics.NewDefaultHTTPServer(ctx)
		if err != nil {
			return err
		}

		return httpServer.ListenAndServe()
	})

	if err := runner.Wait(); err != nil {
		logger.Error().Err(err).Msg("exited with error")
	}
}

// NewHTTPServer creates new HTTP server
func NewHTTPServer(ctx context.Context, cfg *config.Common, chatController api.ServerInterface, _ *openapi3.T) func() error {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(middlewares.UserAuthMiddleware)
	r.Use(middlewares.ResponseLoggerMiddleware)

	api.HandlerFromMux(chatController, r)

	s := &http.Server{
		Handler:           r,
		Addr:              cfg.Server.GetListenPort(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return httputils.HTTPListenAndServeContextErrGroup(ctx, s, 5*time.Second)
}

func initializeDB(ctx context.Context, cfg *config.Common) (conn *sql.DB) {
	conn, err := db.NewPostgres(ctx, cfg.PostgreSQL.BuildDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("connect postgres")
	}
	return conn
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
