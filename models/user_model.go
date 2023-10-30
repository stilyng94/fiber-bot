package models

import (
	"context"
	"time"

	"log/slog"

	"github.com/uptrace/bun"
)

type Role string

const (
	Admin Role = "admin"
	User  Role = "user"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID         int64     `bun:",pk,autoincrement" json:"id"`
	TelegramID int64     `bun:",unique,notnull" json:"telegramId"`
	FirstName  string    `bun:"," json:"firstName,omitempty"`
	LastName   string    `bun:"," json:"lastName,omitempty"`
	Username   string    `bun:"," json:"username,omitempty"`
	CreatedAt  time.Time `bun:",nullzero,notnull,type:timestamp with time zone,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `bun:",nullzero,notnull,type:timestamp with time zone,default:current_timestamp" json:"updatedAt"`
	ChatID     int64     `bun:"," json:"chatId,omitempty"`
	ChatType   string    `bun:"," json:"chatType,omitempty"`
	Role       Role      `bun:",nullzero,notnull,default:'user'" json:"role"`
}

// Create users table.
func CreateUserTable(db *bun.DB, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := db.NewCreateTable().Model((*UserModel)(nil)).IfNotExists().WithForeignKeys().Exec(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Create user table error", slog.Any("error", err))
	}
}
