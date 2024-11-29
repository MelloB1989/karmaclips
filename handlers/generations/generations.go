package generations

import (
	"karmaclips/helpers/generations"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetGenerationsByUserId(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	generations, err := generations.GetGenerationsByUserId(uid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error getting generations",
			"data":    []interface{}{},
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    generations,
	})
}

type GetGenerationsByDateRequest struct {
	Date string `json:"date"`
}

func GetGenerationsByDate(c *fiber.Ctx) error {
	req := new(GetGenerationsByDateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid date format",
		})
	}

	uid := c.Locals("uid").(string)

	generations, err := generations.GetGenerationsByUserIdAndDate(uid, date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error getting generations",
			"data":    []interface{}{},
			"error":   err.Error(),
		})
	}

	if len(generations) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No generations found",
			"data":    []interface{}{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    generations,
	})
}
