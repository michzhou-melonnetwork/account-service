package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"

	"github.com/Melon-Network-Inc/account-service/internal/address"
	"github.com/Melon-Network-Inc/account-service/internal/auth"
	"github.com/Melon-Network-Inc/account-service/internal/config"
	"github.com/Melon-Network-Inc/account-service/internal/errors"
	"github.com/Melon-Network-Inc/account-service/internal/healthcheck"
	"github.com/Melon-Network-Inc/account-service/pkg/accesslog"
	"github.com/Melon-Network-Inc/account-service/pkg/dbcontext"
	"github.com/Melon-Network-Inc/account-service/pkg/log"
	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to the database
	// db, err := dbx.MustOpen("postgres", cfg.DSN)
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), cfg),
	}

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	rg := router.Group("/v1")

	authHandler := auth.Handler(cfg.JWTSigningKey)

	address.RegisterHandlers(rg.Group(""),
		address.NewService(address.NewRepository(db, logger), logger),
		authHandler, logger,
	)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	return router
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
