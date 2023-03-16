package routes

import (
	"kaomojidb/src/handlers/auth"
	"kaomojidb/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/auth")
	// General Middlewares for the route if any

	// Define the subroutes
	route.Get("/login", auth.Login)                     // get your jwt
	route.Get("/renew", middlewares.Auth(), auth.Renew) // Renew your JWT if not blocked
	route.Get("/drop", middlewares.Auth(), auth.Renew)  // drops a given jwt, if none is provided, defaults to the one used.

}
