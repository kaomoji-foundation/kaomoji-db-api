package routes

import (
	"Kaomoji-DB/src/handlers/kaomojis"
	"Kaomoji-DB/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func KaomojisRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/kaomojis")
	// General Middlewares for the route if any

	// Define the subroutes
	route.Post("/", middlewares.Auth(), kaomojis.CreateKaomoji)      // Create
	route.Get("/", middlewares.OptInAuth(), kaomojis.GetKaomojis)    // Read multiple
	route.Get("/:id", middlewares.OptInAuth(), kaomojis.GetKaomoji)  // Read
	route.Patch("/:id", middlewares.Auth(), kaomojis.UpdateKaomoji)  // Update
	route.Delete("/:id", middlewares.Auth(), kaomojis.DeleteKaomoji) // Delete

}
