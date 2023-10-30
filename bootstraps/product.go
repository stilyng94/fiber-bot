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

// BootstrapProduct - Bootstraps the Product router
func BootstrapProduct(router *fiber.App, customValidator *internal.XValidator, sessionStore *session.Store, db *bun.DB, envConfig config.EnvConfig) {
	productRepo := repos.NewProductRepo(db)
	productController := handlers.NewProductHandler(customValidator, productRepo, envConfig)
	productRouter := router.Group("/products", middleware.WithAuth(sessionStore))
	{
		productRouter.Get("/:id", productController.HandleGetProduct)
		productRouter.Get("", productController.HandleGetProducts)
	}

}
