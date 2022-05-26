package kaomojis

import (
	stdMsg "Kaomoji-DB/src/helpers/stdMessages"
	"Kaomoji-DB/src/models"
	"Kaomoji-DB/src/utils/radix"
	"context"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateKaomoji create new kaomoji
// @Summary      Create kaomoji
// @Description  Create a new kaomoji
// @Accept       json
// @Produce      json
// @param	registerData body models.Kaomoji{} true "initial data for the kaomoji"
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      422  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /kaomojis/ [post]
func CreateKaomoji(c *fiber.Ctx) error {
	// return c.SendStatus(fiber.StatusNotImplemented)

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
	if !editorRole.Permissons.KaomojisAdmin {
		return c.Status(fiber.StatusForbidden).JSON(stdMsg.ErrorDefault(
			"failed to update the kaomoji, your role cant update kaomojis.",
			nil,
		))
	}
	// Authenticated & autorized

	/*
	 ———————— Validate fields to be updated, writes changes into kaomoji ———————————
	 Anything not explicitly written here to kaomoji will not be passed to the db

	 							kui >> kaomoji

	*/

	var kaomoji models.Kaomoji
	if err := c.BodyParser(&kaomoji); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Unpocessable Entity, Review your input",
			"data":    err,
		})
	}

	kaomoji.ID = primitive.NilObjectID

	// kaomoji string validation
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

	if kaomoji.Categories == nil || len(kaomoji.Categories) == 0 {
		return c.Status(fiber.StatusConflict).JSON(stdMsg.ErrorDefault(
			"Specified kaomoji string already exists",
			err,
		))
	}
	radix.Sort(kaomoji.Categories)

	// Set the first revision to get the starting point
	kaomoji.Revisions = []models.IssueMinimal{
		{
			ID:        primitive.NewObjectID(),
			Author:    models.UserMinimal{},
			Reviewer:  models.UserMinimal{},
			Title:     "Created by server",
			Status:    "closed",
			Operation: "create",
			Kaomoji:   kaomoji.Minimal(),
		},
	}

	res, err := models.KaomojisCollection.InsertOne(context.Background(), kaomoji)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to create the kaomoji",
			"data":    nil,
		})
	}

	err = kaomoji.Fill(res.InsertedID.(primitive.ObjectID).Hex(), true, false)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to retrieve the data of the created resource 'kaomoji', with id: " + res.InsertedID.(primitive.ObjectID).Hex(),
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Kaomoji successfully created",
		"data":    kaomoji,
	})
}
