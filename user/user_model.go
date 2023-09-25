package user

import (
	"context"
	"time"

	"log/slog"

	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID         int64     `bun:",pk,autoincrement" json:"id"`
	Name       string    `bun:",notnull" json:"name"`
	Email      string    `bun:",notnull,unique" json:"email"`
	Password   string    `bun:",notnull" json:"-"`
	IsVerified bool      `bun:",notnull,default:0" json:"isVerified"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:now" json:"createdAt"`
	UpdatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

// Create users table.
func CreateTable(db *bun.DB, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := db.NewCreateTable().Model((*UserModel)(nil)).IfNotExists().WithForeignKeys().Exec(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Create user table error", slog.Any("error", err))
	}
}

// Drop and create tables.
func ResetTable(db *bun.DB, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	err := db.ResetModel(ctx, (*UserModel)(nil))
	if err != nil {
		logger.ErrorContext(ctx, "Reset user table error", slog.Any("error", err))
	}
}
