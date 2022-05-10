package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(api *fiber.Router) {
	// Middleware if any
	//? wouldn't this section be useless?, since it would be the same as global middlewares

	// Routes /api/v1
	//(*api).Get("/docs/*", swagger.New()) // mount the docs

	DocsRoute(api) // >> /docs mount the docs trough files

	UsersRoute(api) // >> /users
	AuthRoute(api)  // >> /auth

	PingRoute(api) // >> /ping

}
