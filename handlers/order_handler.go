package handlers

import (
	"encoding/json"
	"strconv"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/fiber-bot/config"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/models"
	"github.com/fiber-bot/repos"
	"github.com/fiber-bot/tele_bot"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type OrderHandler struct {
	customValidator *internal.XValidator
	sessionStore    *session.Store
	OrderRepo       repos.OrderRepo
	envConfig       config.EnvConfig
}

func NewOrderHandler(customValidator *internal.XValidator, sessionStore *session.Store, orderRepo repos.OrderRepo, envConfig config.EnvConfig) *OrderHandler {
	return &OrderHandler{
		customValidator: customValidator,
		sessionStore:    sessionStore,
		OrderRepo:       orderRepo,
		envConfig:       envConfig,
	}
}

// [HandleGetOrders] responds with the list of all Orders as JSON.
func (handler *OrderHandler) HandleGetOrders(c *fiber.Ctx) error {
	orders, err := handler.OrderRepo.GetOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": orders})
}

// [HandleGetOrder] responds with the Order as JSON.
func (handler *OrderHandler) HandleGetOrder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	order, err := handler.OrderRepo.GetOrder(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": order})
}

// [HandleAddToCart] responds with the Order as JSON.
func (handler *OrderHandler) HandleAddToCart(c *fiber.Ctx) error {

	var payload struct {
		Items []models.CartItem `json:"items" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	if errors := handler.customValidator.Validate(payload); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": errors})
	}

	cartByte, err := json.Marshal(payload.Items)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	sess, _ := handler.sessionStore.Get(c)
	sess.Set("cart", string(cartByte))
	_ = sess.Save()

	return c.JSON(fiber.Map{"message": "success"})
}

// [HandleCheckout]
func (handler *OrderHandler) HandleCheckout(c *fiber.Ctx) error {
	sess, _ := handler.sessionStore.Get(c)
	userID := sess.Get("userId")
	cartSess := sess.Get("cart")
	if cartSess == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cart is empty"})
	}

	var cartItems []models.CartItem
	err := json.Unmarshal([]byte(cartSess.(string)), &cartItems)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var orderItems []models.OrderItem
	var telePrices []gotgbot.LabeledPrice
	var totalCharge int64 = 0
	for _, item := range cartItems {
		itemTotalPrice := 10 * item.Quantity
		totalCharge += itemTotalPrice
		telePrices = append(telePrices, gotgbot.LabeledPrice{Label: item.Title, Amount: itemTotalPrice})
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Title:     item.Title,
			Price:     itemTotalPrice,
		})
	}

	b, err := tele_bot.ConfigBot(handler.envConfig, false)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order, err := handler.OrderRepo.InsertOrder(c.Context(), repos.HandleAddOrderPayload{
		UserID:     userID.(int64),
		Amount:     totalCharge,
		OrderItems: orderItems,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	inv, err := b.CreateInvoiceLink(
		"awesome shop",
		"Payment for items",
		strconv.FormatInt(order.ID, 10),
		handler.envConfig.TelegramStripeToken,
		"USD",
		telePrices, nil,
	)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"invoice": inv, "orderId": order.ID, "totalCharge": totalCharge})
}

// [HandleOrderUpdate]
func (handler *OrderHandler) HandleOrderUpdate(c *fiber.Ctx) error {
	sess, _ := handler.sessionStore.Get(c)
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	var payload struct {
		Status models.OrderStatus `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = handler.OrderRepo.UpdateOrder(c.Context(), id, payload.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sess.Delete("cart")
	return c.JSON(fiber.Map{"message": "success"})
}
