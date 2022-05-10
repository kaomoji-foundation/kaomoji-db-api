package routes

import (
	"GO-API-template/src/handlers/auth"
	"GO-API-template/src/middlewares"

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
