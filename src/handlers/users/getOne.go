package users

import (
	"GO-API-template/src/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUser get a user
// @Summary      Retrieve user data
// @Description  Check api is active
// @security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uid  query  string  true  "uid string"
// @Router       /users/{uid} [get]
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      404  {object}  interface{}
// @Failure      500  {object}  interface{}
func GetUser(c *fiber.Ctx) error {
	// Identity of the user to get data from
	identity := c.Params("uid")
	// get the data of the user we want to get data from
	var user models.User
	err := user.Fill(identity, true, true, false)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with ID", "data": nil})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Found an error while trying to get the user",
		})
	}

	// Use public profile if the user is not authenticated
	if !c.Locals("Authenticated").(bool) {
		return c.JSON(fiber.Map{"status": "success", "message": "User found", "user": user.Public()})
	}

	// Get Token of the reader's user
	token := c.Locals("user").(*jwt.Token)
	editorUID := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["uid"])

	// Get the reader's data
	var editorUser models.User
	editorUser.Fill(editorUID, true, false, false)
	var editorRole models.Role
	editorRole.Fill(editorUser.RoleID.Hex(), true, false)

	// check what the user is authorised to get and return that
	if user.ID.Hex() == editorUID {
		// User owner
		return c.JSON(fiber.Map{"status": "success", "message": "User found", "user": user.Private()})
	}
	// Parametrized permissons
	if !editorRole.Permissons.UsersAdmin {
		// Any regular user:
		return c.JSON(fiber.Map{"status": "success", "message": "User found", "user": user.Public()})
	}
	// User admins
	return c.JSON(fiber.Map{"status": "success", "message": "User found", "user": user})
}
