package handlers

import (
	"strconv"

	"github.com/fiber-bot/config"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/repos"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UserHandler struct {
	customValidator *internal.XValidator
	sessionStore    *session.Store
	userRepo        repos.UserRepo
	envConfig       config.EnvConfig
}

func NewUserHandler(customValidator *internal.XValidator, sessionStore *session.Store, userRepo repos.UserRepo, envConfig config.EnvConfig) *UserHandler {
	return &UserHandler{
		customValidator: customValidator,
		sessionStore:    sessionStore,
		userRepo:        userRepo,
		envConfig:       envConfig,
	}
}

// [HandleGetUsers] responds with the list of all users as JSON.
func (handler *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := handler.userRepo.GetUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": users})
}

// [HandleGetUser] responds with the user as JSON.
func (handler *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := handler.userRepo.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": user})
}

// [HandleUpsertUser]
func (handler *UserHandler) HandleUpsertUser(c *fiber.Ctx) error {
	var payload struct {
		Token string `json:"token" validate:"required"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	if errors := handler.customValidator.Validate(payload); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": errors})
	}
	parsedUser, err := internal.DecodeTelegramHash(payload.Token, handler.envConfig)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := handler.userRepo.CreateUser(c.Context(), parsedUser)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	sess, _ := handler.sessionStore.Get(c)
	_ = sess.Regenerate()
	sess.Set("queryId", parsedUser.QueryID)
	sess.Set("isAuthenticated", 1)
	sess.Set("userId", user.ID)
	_ = sess.Save()
	return c.JSON(fiber.Map{"data": user})
}
