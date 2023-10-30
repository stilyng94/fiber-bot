package bootstraps

import (
	"github.com/fiber-bot/config"
	"github.com/fiber-bot/handlers"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/repos"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/uptrace/bun"
)

// BootstrapUser - initializes the user router
func BootstrapUser(router *fiber.App, customValidator *internal.XValidator, sessionStore *session.Store, db *bun.DB, envConfig config.EnvConfig) {
	userRepo := repos.NewUserRepo(db)
	userController := handlers.NewUserHandler(customValidator, sessionStore, userRepo, envConfig)
	userRouter := router.Group("/users")
	{
		userRouter.Get("/:id", userController.HandleGetUser)
		userRouter.Get("", userController.HandleGetUsers)
		userRouter.Post("", userController.HandleUpsertUser)
	}

}
