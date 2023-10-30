package repos

import (
	"context"
	"encoding/json"

	"github.com/fiber-bot/models"
	"github.com/uptrace/bun"
)

type HandleAddOrderPayload struct {
	UserID     int64              `json:"userId" validate:"required"`
	Amount     int64              `json:"amount" validate:"required,gt=0"`
	OrderItems []models.OrderItem `json:"orderItems" validate:"required"`
}

type OrderRepo interface {
	InsertOrder(ctx context.Context, payload HandleAddOrderPayload) (*models.OrderModel, error)
	UpdateOrder(ctx context.Context, ID int, status models.OrderStatus) (*models.OrderModel, error)
	DeleteOrder(ctx context.Context, ID int) error
	GetOrder(ctx context.Context, ID int) (*models.ExpandedOrderModel, error)
	GetOrders(ctx context.Context) ([]models.ExpandedOrderModel, error)
}

type OrderRepoImpl struct {
	db *bun.DB
}

// GetOrders implements OrderRepo.
func (u *OrderRepoImpl) GetOrders(ctx context.Context) ([]models.ExpandedOrderModel, error) {
	orders := []models.OrderModel{}
	err := u.db.NewSelect().Model(&orders).Limit(10).Scan(ctx)
	if err != nil {
		return nil, err
	}
	var expandedOrders []models.ExpandedOrderModel
	for _, v := range orders {
		var expandedOrder models.ExpandedOrderModel
		err = json.Unmarshal([]byte(v.OrderItems), &expandedOrder)
		if err != nil {
			return nil, err
		}
		expandedOrder.ID = v.ID
		expandedOrder.CreatedAt = v.CreatedAt
		expandedOrder.Status = v.Status
		expandedOrder.User = v.User
		expandedOrder.UserID = v.UserID
		expandedOrders = append(expandedOrders, expandedOrder)
	}

	return expandedOrders, nil
}

// DeleteOrder implements OrderRepo.
func (u *OrderRepoImpl) DeleteOrder(ctx context.Context, ID int) error {
	Order := &models.OrderModel{ID: int64(ID)}
	_, err := u.db.NewDelete().Model(Order).WherePK().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// GetOrder implements OrderRepo.
func (u *OrderRepoImpl) GetOrder(ctx context.Context, ID int) (*models.ExpandedOrderModel, error) {
	var order models.OrderModel
	err := u.db.NewSelect().Model(&order).Where("id = ?", int64(ID)).Scan(ctx)
	if err != nil {
		return nil, err
	}
	var expandedOrder models.ExpandedOrderModel
	err = json.Unmarshal([]byte(order.OrderItems), &expandedOrder)
	if err != nil {
		return nil, err
	}
	expandedOrder.ID = order.ID
	expandedOrder.CreatedAt = order.CreatedAt
	expandedOrder.Status = order.Status
	expandedOrder.User = order.User
	expandedOrder.UserID = order.UserID

	return &expandedOrder, nil
}

// UpdateOrder implements OrderRepo.
func (u *OrderRepoImpl) UpdateOrder(ctx context.Context, ID int, status models.OrderStatus) (*models.OrderModel, error) {
	Order := &models.OrderModel{ID: int64(ID), Status: status}
	_, err := u.db.NewUpdate().Model(Order).Column("status").WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}
	return Order, nil
}

// InsertOrder implements OrderRepo.
func (u *OrderRepoImpl) InsertOrder(ctx context.Context, payload HandleAddOrderPayload) (*models.OrderModel, error) {
	orderItems, err := json.Marshal(payload.OrderItems)
	if err != nil {
		return nil, err
	}
	order := &models.OrderModel{UserID: payload.UserID, Amount: payload.Amount,
		OrderItems: string(orderItems)}
	_, err = u.db.NewInsert().Model(order).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func NewOrderRepo(db *bun.DB) OrderRepo {
	return &OrderRepoImpl{db: db}
}
