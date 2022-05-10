package routes

import (
	"github.com/gofiber/fiber/v2"
)

// Ping is a route to check if server is woken
// @Summary      API docs
// @Description  get the API docs, in json, yaml, or view them using redoc in /docs/spec.html
// @security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  string    string
// @Failure      401  {object}  interface{}
// @Router       /docs [get]
func DocsRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/docs")
	// General Middlewares for the route if any
	route.Get("/", func(c *fiber.Ctx) error { // redirect to the web docs
		return c.Redirect("docs/spec.htm")
	})
	// Define the subroutes
	route.Static("/", "./docs/public") // mount the docs

}
