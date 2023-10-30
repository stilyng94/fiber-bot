package models

import (
	"context"
	"time"

	"log/slog"

	"github.com/uptrace/bun"
)

type OrderStatus string

const (
	OrderStatusSuccess   OrderStatus = "paid"
	OrderStatusFailure   OrderStatus = "failed"
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type CartItem struct {
	ProductID int64  `json:"productId" validate:"required"`
	Quantity  int64  `json:"quantity" validate:"required"`
	Title     string `json:"title" validate:"required"`
}

type OrderItem struct {
	ProductID int64  `json:"productId"`
	Quantity  int64  `json:"quantity"`
	Price     int64  `json:"price"`
	Title     string `json:"title"`
}

type OrderModel struct {
	bun.BaseModel `bun:"table:orders,alias:o"`

	ID         int64       `bun:",pk,autoincrement"`
	UserID     int64       `bun:",notnull"`
	Amount     int64       `bun:",notnull"`
	OrderItems string      `bun:",notnull" json:"orderItems"`
	Status     OrderStatus `bun:",nullzero,notnull,default:'pending'"`
	User       *UserModel  `bun:",rel:belongs-to,join_on:user_id=id" json:"omitempty"`
	CreatedAt  time.Time   `bun:",nullzero,notnull,type:timestamp with time zone,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time   `bun:",nullzero,notnull,type:timestamp with time zone,default:current_timestamp" json:"updatedAt"`
}

type ExpandedOrderModel struct {
	ID         int64
	UserID     int64
	Amount     int64
	OrderItems []OrderItem `json:"orderItems"`
	Status     OrderStatus
	User       *UserModel `json:"omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}

// Create orders table.
func CreateOrderTable(db *bun.DB, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := db.NewCreateTable().Model((*OrderModel)(nil)).IfNotExists().WithForeignKeys().Exec(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Create order table error", slog.Any("error", err))
	}
}
