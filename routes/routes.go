package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"karmaclips/database"
	"karmaclips/handlers/auth"
	"karmaclips/handlers/gen"
	"karmaclips/handlers/generations"
	"karmaclips/handlers/services"
	"karmaclips/middlewares"
)

func Routes() *fiber.App {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, X-Karma-Admin-Auth, Authorization",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(database.ResponseHTTP{
			Success: true,
			Data:    nil,
			Message: "OK",
		})
	})
	v1 := app.Group("/v1")
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(database.ResponseHTTP{
			Success: true,
			Data:    nil,
			Message: "OK",
		})
	})

	authRoutes := v1.Group("/auth")
	authRoutes.Post("/login", auth.Login)
	authRoutes.Post("/register", auth.Register)

	genRoutes := v1.Group("/gen")
	genRoutes.Post("/image", middlewares.IsUserVerified, gen.CreateImage)
	genRoutes.Get("/job/:jobId", middlewares.IsUserVerified, gen.GetJobStatus)

	generationRoutes := v1.Group("/generations")
	generationRoutes.Get("/", middlewares.IsUserVerified, generations.GetGenerationsByUserId)
	generationRoutes.Post("/bydate", middlewares.IsUserVerified, generations.GetGenerationsByDate)

	servicesRoutes := v1.Group("/services")
	servicesRoutes.Post("/", middlewares.IsUserVerified, services.CreateService)
	servicesRoutes.Get("/", middlewares.IsUserVerified, services.GetServices)
	servicesRoutes.Get("/:id", middlewares.IsUserVerified, services.GetServiceById)

	return app
}
