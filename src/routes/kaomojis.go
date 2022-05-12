package routes

import (
	"GO-API-template/src/handlers/kaomojis"
	"GO-API-template/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func KaomojisRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/kaomojis") //? shuld i use /user instead? it makes it a lot more semantic
	// General Middlewares for the route if any

	// Define the subroutes
	route.Post("/", kaomojis.CreateKaomoji)                          // Create
	route.Get("/", kaomojis.GetKaomojis)                             // Read multiple
	route.Get("/:id", middlewares.OptInAuth(), kaomojis.GetKaomoji)  // Read
	route.Patch("/:id", middlewares.Auth(), kaomojis.UpdateKaomoji)  // Update
	route.Delete("/:id", middlewares.Auth(), kaomojis.DeleteKaomoji) // Delete

}
