package generations

import (
	"karmaclips/helpers/generations"

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
