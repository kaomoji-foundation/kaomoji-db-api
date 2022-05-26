package kaomojis

import (
	stdMsg "Kaomoji-DB/src/helpers/stdMessages"
	"Kaomoji-DB/src/models"
	"Kaomoji-DB/src/utils/radix"
	"context"
	"reflect"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// UpdateKaomoji update kaomoji
// @Summary      update kaomoji
// @Description  Update kaomoji info
// @Accept       json
// @Produce      json
// @security     BearerAuth
// @param	updateKaomojiData body models.Kaomoji{} true "data to update, currently only allows to update the fullName field"
// @param	id path string true "Kaomoji ID or kaomoji string"
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      422  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /kaomojis/{id} [patch]
func UpdateKaomoji(c *fiber.Ctx) error {
	// kaomoji update input
	var kui models.Kaomoji
	if err := c.BodyParser(&kui); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Unpocessable Entity, Review your input",
			"data":    err,
		})
	}

	// Token of the editor's user
	token := c.Locals("user").(*jwt.Token)
	editorUID := token.Claims.(jwt.MapClaims)["uid"].(string)

	// Get the editor's data
	var editorUser models.User
	err := editorUser.Fill(editorUID, true, false, false)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the editor's user with ID: "+editorUID,
				err.Error(),
			))
	}
	var editorRole models.Role
	err = editorRole.Fill(editorUser.RoleID.Hex(), true, false)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the editor's role with ID: "+editorUser.RoleID.Hex(),
				err.Error(),
			))
	}

	// check if editor is authorised to do the operation using Parametrized permissons
	if !editorRole.Permissons.KaomojisModerator || !editorRole.Permissons.KaomojisAdmin {
		return c.Status(fiber.StatusForbidden).JSON(stdMsg.ErrorDefault(
			"failed to update the kaomoji, your role cant update kaomojis.",
			nil,
		))
	}
	// Authenticated & autorized

	// Identity of the kaomoji to modify
	identity := c.Params("id")
	// get the data of the kaomoji we want to modify
	var kaomoji models.Kaomoji
	err = kaomoji.Fill(identity, true, true)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the requested user with identificator: "+identity,
				err.Error(),
			))
	}

	/*
	 ———————— Validate fields to be updated, writes changes into kaomoji ———————————
	 Anything not explicitly written here to kaomoji will not be passed to the db

	 							kui >> kaomoji

	*/

	// Editing revisions requires special perms, adding doesn't
	if len(kui.Revisions) != len(kaomoji.Revisions) {
		if len(kui.Revisions) > len(kaomoji.Revisions) {
			kaomoji.Revisions = kui.Revisions
		} else if editorRole.Permissons.KaomojisAdmin {
			kaomoji.Revisions = kui.Revisions
		} else {
			return c.Status(fiber.StatusForbidden).JSON(stdMsg.ErrorDefault(
				"failed to update the kaomoji, your role cant update kaomojis.",
				nil,
			))
		}
	}

	// kaomoji string validation
	if kui.String != "" || len(kui.String) > 0 || !reflect.DeepEqual(kaomoji.String, kui.String) {
		kaomoji.String = kui.String

		unique, err := kaomoji.CheckUnique()
		if err != nil {
			return c.Status(500).JSON(stdMsg.ErrorDefault(
				"An error ocured while checking if kaomoji string is unique",
				err,
			))
		}
		if !unique {
			return c.Status(fiber.StatusConflict).JSON(stdMsg.ErrorDefault(
				"Specified kaomoji string already exists",
				err,
			))
		}
	}
	if kui.Desciption != "" || kui.Desciption != kaomoji.Desciption {
		kaomoji.Desciption = kui.Desciption
	}

	if kui.Categories != nil || len(kui.Categories) > 0 {
		radix.Sort(kui.Categories)
		kaomoji.Categories = kui.Categories
	}

	filter := bson.M{"_id": kaomoji.ID}
	update := bson.M{"$set": kaomoji}

	_, err = models.KaomojisCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to update the kaomoji",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Kaomoji successfully updated",
		"data":    kaomoji,
	})
}
