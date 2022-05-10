package routes

import (
	"GO-API-template/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

// Ping is a route to check if server is woken
// @Summary      Ping the server
// @Description  Check api is active
// @security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  string    string
// @Failure      401  {object}  interface{}
// @Router       /ping [get]
func PingRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/ping")
	// General Middlewares for the route if any

	// Define the subroutes
	route.Get("/", middlewares.Auth(), func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Ping!")
	})
}
