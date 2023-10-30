package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func WithAuth(sessionStore *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, _ := sessionStore.Get(c)
		userID := sess.Get("userId")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthenticated"})
		}
		isAuthenticated := sess.Get("isAuthenticated")
		if isAuthenticated == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthenticated"})
		}
		return c.Next()
	}
}
