package services

import (
	"karmaclips/database"
	"karmaclips/helpers/services"

	"github.com/gofiber/fiber/v2"
)

func CreateService(c *fiber.Ctx) error {
	var req interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	switch req.(type) {
	case map[string]interface{}:
		// Single service
		serviceReq := new(database.AiServices)
		if err := c.BodyParser(serviceReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request",
			})
		}

		service, err := services.CreateService(serviceReq)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Service creation failed",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Service created successfully",
			"data":    service,
		})

	case []interface{}:
		// Array of services
		var servicesReq []database.AiServices
		if err := c.BodyParser(&servicesReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request",
			})
		}

		var createdServices []database.AiServices
		for _, serviceReq := range servicesReq {
			service, err := services.CreateService(&serviceReq)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Service creation failed",
				})
			}
			createdServices = append(createdServices, *service)
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Services created successfully",
			"data":    createdServices,
		})

	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request format",
		})
	}
}

func GetServiceById(c *fiber.Ctx) error {
	sid := c.Params("sid")

	service, err := services.GetServiceById(sid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "No service found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Service found",
		"data":    service,
	})
}

func GetServices(c *fiber.Ctx) error {
	services, err := services.GetServices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "No services found",
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Services found",
		"data":    services,
	})
}
