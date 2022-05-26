package routes

import (
	"Kaomoji-DB/src/handlers/auth"
	"Kaomoji-DB/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuthRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/auth")
	// General Middlewares for the route if any

	// Define the subroutes
	route.Get("/login", auth.Login)                     // get your jwt
	route.Get("/renew", middlewares.Auth(), auth.Renew) // Renew your JWT if not blocked

}
