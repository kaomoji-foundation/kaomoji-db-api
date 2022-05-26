package kaomojis

import (
	"Kaomoji-DB/src/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetKaomoji get a kaomoji
// @Summary      Retrieve kaomoji data
// @Description  Check api is active
// @security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        uid  query  string  true  "id string"
// @Router       /kaomojis/{id} [get]
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      404  {object}  interface{}
// @Failure      500  {object}  interface{}
func GetKaomoji(c *fiber.Ctx) error {
	// Identificator of the kaomoji to get data from
	identity := c.Params("id")
	// get the data of the kaomoji we want to get data from
	var kaomoji models.Kaomoji
	err := kaomoji.Fill(identity, true, true)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No kaomoji found with ID", "data": nil})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Found an error while trying to get the kaomoji",
		})
	}

	// Use minimal version if the user is not authenticated to save bandwith
	if !c.Locals("Authenticated").(bool) {
		return c.JSON(fiber.Map{"status": "success", "message": "Kaomoji found", "kaomoji": kaomoji.Minimal()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Kaomoji found", "kaomoji": kaomoji})
}
