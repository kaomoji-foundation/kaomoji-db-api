package kaomojis

import (
	stdMsg "Kaomoji-DB/src/helpers/stdMessages"
	"Kaomoji-DB/src/models"
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// DeleteKaomoji delete kaomoji
// @Summary      delete kaomoji
// @Description  delete kaomoji completely
// @Accept       json
// @Produce      json
// @security     BearerAuth
// @param	password body PasswordInput{} false "password of the kaomoji to delete, not required if kaomoji is admin"
// @param	uid path string true "Kaomoji ID"
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      422  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /kaomojis/{uid} [delete]
func DeleteKaomoji(c *fiber.Ctx) error {

	// Token of the editor's user
	token := c.Locals("user").(*jwt.Token)
	editorUID := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["uid"])

	// Get the editor's data
	var editorUser models.User
	editorUser.Fill(editorUID, true, false, false)
	var editorRole models.Role
	editorRole.Fill(editorUser.RoleID.Hex(), true, false)

	// Parametrized permissons
	if !editorRole.Permissons.KaomojisAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete the kaomoji, your role cant delete other kaomojis.",
			"data":    nil,
		})
	}

	// Authenticated & autorized

	// Identity of the kaomoji to modify
	kaomojiIdentificator := c.Params("uid")
	var kaomoji models.Kaomoji
	err := kaomoji.Fill(kaomojiIdentificator, true, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(stdMsg.ErrorDefault("Something went wrong while fetching the kaomoji drom the DB", nil))
	}

	filter := bson.M{"_id": kaomoji.ID}
	deleted, err := models.KaomojisCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something went wron while deleting the kaomoji", "data": nil})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("%v kaomoji successfully deleted", deleted.DeletedCount),
		"data":    kaomoji,
	})
}
