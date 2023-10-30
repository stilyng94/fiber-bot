package config

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// InitDB - Connect to database and run migrations
func InitDB(envConfig EnvConfig) (*bun.DB, error) {
	drv, err := connectToDriver(envConfig)
	if err != nil {
		return nil, err
	}

	var client *bun.DB
	if envConfig.DatabaseEngine == postgresEngine {
		client = bun.NewDB(drv, pgdialect.New())
	} else {
		client = bun.NewDB(drv, sqlitedialect.New())
	}

	client.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(envConfig.IsDev()),
		bundebug.WithEnabled(envConfig.IsDev()),
	))
	return client, nil
}

// PingDB - check if database is healthy
func PingDB(db *bun.DB) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := db.NewSelect().ColumnExpr("1").Exec(ctx)
	return err != nil
}

func connectToDriver(envConfig EnvConfig) (*sql.DB, error) {
	if envConfig.DatabaseEngine == postgresEngine {
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(envConfig.DSN)))
		if sqldb == nil {
			return nil, errors.New("database connection error")
		}
		return sqldb, nil
	}

	sqldb, err := sql.Open(sqliteshim.ShimName, envConfig.DSN)
	if err != nil {
		return nil, err
	}
	return sqldb, nil
}

// Drop and create tables.
func ResetTable(db *bun.DB, logger *slog.Logger, models ...interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err := db.ResetModel(ctx, models...)
	if err != nil {
		logger.ErrorContext(ctx, "Reset tables error", slog.Any("error", err))
	}
}
