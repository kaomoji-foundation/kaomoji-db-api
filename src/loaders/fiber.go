package loaders

import (
	"GO-API-template/src/middlewares"
	"GO-API-template/src/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// This will load the routes onto a fiber app, configured here as well
func LoadFiber() *fiber.App {
	cfg := fiber.Config{
		CaseSensitive:     true,
		EnablePrintRoutes: false,
	}
	app := fiber.New(cfg)

	//* here is where middlewares used in all routes should be mounted
	app.Use(middlewares.OpenTelemery()) // use open-telemetry middleware with custom config
	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New())
	// limmit to 5 requests per 10 seconds max
	app.Use(limiter.New(limiter.Config{
		Expiration: 30,
		Max:        20,
	}))
	//* here you mount the routes for the apps in a certain bas path like: /api/v1
	router := app.Group("/v1")  // this is the base route for all endpoints
	routes.SetupRoutes(&router) // this Mounts all the app routes into router

	return app
}
