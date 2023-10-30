package models

import (
	"context"
	"log/slog"
	"time"

	"github.com/uptrace/bun"
)

type ProductStatus string

const (
	ProductStatusAvailable   ProductStatus = "available"
	ProductStatusUnavailable ProductStatus = "unavailable"
)

type ProductModel struct {
	bun.BaseModel `bun:"table:products,alias:p"`

	ID          int64         `bun:",pk,autoincrement"`
	Title       string        `bun:","`
	ImageUrl    string        `bun:","`
	Price       int64         `bun:","`
	Description string        `bun:","`
	Status      ProductStatus `bun:",nullzero,notnull,default:'available'"`
	CreatedAt   time.Time     `bun:",nullzero,notnull,type:timestamp with time zone,default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time     `bun:",nullzero,notnull,type:timestamp with time zone,default:current_timestamp" json:"updatedAt"`
}

// Create product table.
func CreateProductTable(db *bun.DB, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := db.NewCreateTable().Model((*ProductModel)(nil)).IfNotExists().WithForeignKeys().Exec(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Create product table error", slog.Any("error", err))
	}
}
