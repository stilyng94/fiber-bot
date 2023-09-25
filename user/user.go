package user

import (
	"github.com/fiber-bot/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/uptrace/bun"
)

// InitializeUser - initializes the user router
func InitializeUser(router *fiber.App, customValidator *internal.XValidator, sessionStore *session.Store, db *bun.DB) {
	userRepo := NewUserRepo(db)
	userController := NewUserHandler(customValidator, sessionStore, userRepo)
	userRouter := router.Group("/users")
	{
		userRouter.Get("/:id", userController.HandleGetUser)
		userRouter.Get("", userController.HandleGetUsers)
		userRouter.Post("", userController.HandleAddUser)

	}

}
