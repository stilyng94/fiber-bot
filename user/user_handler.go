package user

import (
	"strconv"

	"github.com/fiber-bot/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type UserHandler struct {
	customValidator *internal.XValidator
	sessionStore    *session.Store
	userRepo        UserRepo
}

type HandleAddUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func NewUserHandler(customValidator *internal.XValidator, sessionStore *session.Store, userRepo UserRepo) *UserHandler {
	return &UserHandler{
		customValidator: customValidator,
		sessionStore:    sessionStore,
		userRepo:        userRepo,
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

// [HandleAddUser] responds with the User as JSON.
func (handler *UserHandler) HandleAddUser(c *fiber.Ctx) error {
	var payload HandleAddUserPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	if errors := handler.customValidator.Validate(payload); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": errors})
	}

	user, err := handler.userRepo.InsertUser(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": user})
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
