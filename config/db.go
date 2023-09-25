package config

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// InitDB - Connect to database and run migrations
func InitDB(envConfig EnvConfig) (*bun.DB, error) {
	drv, err := connectToDriver(envConfig.DSN)
	if err != nil {
		return nil, err
	}
	client := bun.NewDB(drv, sqlitedialect.New())
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

func connectToDriver(dsn string) (*sql.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		return nil, err
	}
	return sqldb, nil
}
