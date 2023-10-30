package handlers

import (
	"strconv"

	"github.com/fiber-bot/config"
	"github.com/fiber-bot/internal"
	"github.com/fiber-bot/repos"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	customValidator *internal.XValidator
	ProductRepo     repos.ProductRepo
	envConfig       config.EnvConfig
}

func NewProductHandler(customValidator *internal.XValidator, orderRepo repos.ProductRepo, envConfig config.EnvConfig) *ProductHandler {
	return &ProductHandler{
		customValidator: customValidator,
		ProductRepo:     orderRepo,
		envConfig:       envConfig,
	}
}

// [HandleGetProducts] responds with the list of all Products as JSON.
func (handler *ProductHandler) HandleGetProducts(c *fiber.Ctx) error {
	products, err := handler.ProductRepo.GetProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": products})
}

// [HandleGetProduct] responds with the Product as JSON.
func (handler *ProductHandler) HandleGetProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}
	order, err := handler.ProductRepo.GetProduct(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": order})
}
