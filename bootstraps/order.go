package bootstraps

import (
	"github.com/fiber-bot/config"
	"github.com/fiber-bot/handlers"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/middleware"
	"github.com/fiber-bot/repos"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/uptrace/bun"
)

// BootstrapOrder - Bootstraps the Order router
func BootstrapOrder(router *fiber.App, customValidator *internal.XValidator, sessionStore *session.Store, db *bun.DB, envConfig config.EnvConfig) {
	orderRepo := repos.NewOrderRepo(db)
	orderController := handlers.NewOrderHandler(customValidator, sessionStore, orderRepo, envConfig)
	orderRouter := router.Group("/orders", middleware.WithAuth(sessionStore))
	{
		orderRouter.Post("/cart", orderController.HandleAddToCart)
		orderRouter.Get("/:id", orderController.HandleGetOrder)
		orderRouter.Get("", orderController.HandleGetOrders)
		orderRouter.Post("", orderController.HandleCheckout)
	}

}
