package routes

import (
	"GO-API-template/src/handlers/users"
	"GO-API-template/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UsersRoute(r *fiber.Router) {
	// Start the route
	route := (*r).Group("/users") //? shuld i use /user instead? it makes it a lot more semantic
	// General Middlewares for the route if any

	// Define the subroutes
	route.Post("/", users.CreateUser)                           // Create
	route.Get("/", users.GetUsers)                              // Create
	route.Get("/:uid", middlewares.OptInAuth(), users.GetUser)  // Read
	route.Patch("/:uid", middlewares.Auth(), users.UpdateUser)  // Update
	route.Delete("/:uid", middlewares.Auth(), users.DeleteUser) // Delete

}
