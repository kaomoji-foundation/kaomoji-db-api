package users

import (
	stdMsg "GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"GO-API-template/src/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser register new user
// @Summary      Register endpoint
// @Description  Register a new user
// @Accept       json
// @Produce      json
// @param	registerData body models.User{} true "initial data for the user"
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      422  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /users/ [post]
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Unpocessable Entity, Review your input",
			"data":    err,
			"required": fiber.Map{
				"username": "SuperCoolUsername",
				"email":    "you@example.com",
				"password": "YourSuperSecurePassword",
				"fullName": "yourName",
			},
		})
	}

	// check that the username/email is not already being used
	unique, err := user.CheckUnique()
	if err != nil {
		return c.Status(500).JSON(stdMsg.ErrorDefault(
			"An error ocured while checking if username and email are unique",
			err,
		))
	}
	if !unique {
		return c.Status(fiber.StatusConflict).JSON(stdMsg.ErrorDefault(
			"Specified username or email is already being used",
			err,
		))
	}

	hash, err := utils.Hash(user.Password)
	if err != nil {
		return c.Status(500).JSON(stdMsg.ErrorDefault("Couldn't hash password", err))

	}

	user.Password = hash

	// Lock to create onlyregular user by asignin the "user" role
	user.Role = "user"
	// Set the roleID dynamically so it can be para metrized in the future
	err = user.SetRole()

	// lock tokens to be empty at user creation
	user.Tokens = map[string]bool{}
	user.BlockedTokens = map[string]bool{}

	// lock id to force generating a new one by the db
	user.ID = primitive.NilObjectID

	_, err = user.Create()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Created user", "data": user.Private()})
}
