package users

import (
	"GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"context"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rangeUser struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Total   int64       `json:"total"`
	Count   int         `json:"count"`
	Offset  int         `json:"offset"`
	Limmit  int         `json:"limmit"`
	Next    string      `json:"next"`
	Users   interface{} `json:"users"`
}

// GetUsers get the users list
// @Summary      Retrieve users list
// @Description  Retrieve the users id's list
// @security     BearerAuth
// @Accept       json
// @Produce      json
// @Router       /users [get]
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      404  {object}  interface{}
// @Failure      500  {object}  interface{}
func GetUsers(c *fiber.Ctx) error {
	offset, offsetErr := strconv.Atoi(c.Query("o", "0"))
	limmit, limmitErr := strconv.Atoi(c.Query("l", "10"))
	if (limmit - offset) > 100 {
		limmit = offset + 100
	}
	if offsetErr != nil || limmitErr != nil {
		offset = 0
		limmit = 10
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Found an error while parsing your input, review your 'o' and 'l' query params",
		})
	}

	projection := bson.M{"_id": 1, "username": 1, "role": 1} // to match models.user.minimal
	cursor, err := models.UsersCollection.Find(
		context.Background(),
		bson.D{},
		options.Find().SetSkip(int64(offset)).SetLimit(int64(limmit)).SetProjection(projection))
	if err != nil && err != mongo.ErrNoDocuments {
		c.Status(fiber.StatusInternalServerError).JSON(
			stdMessages.ErrorDefault("An error ocurred while retrieving the users data", nil),
		)
	}

	var users []models.UserMinimal
	err = cursor.All(context.Background(), &users)
	if err != nil && err != mongo.ErrNoDocuments {
		c.Status(fiber.StatusInternalServerError).JSON(
			stdMessages.ErrorDefault("An error ocurred while retrieving the users data", nil),
		)
	}

	r := limmit - offset
	var next string
	if r <= len(users) {
		next = c.BaseURL() + string(c.Request().URI().Path()) + fmt.Sprintf("/?o=%v&l=%v", offset+r, limmit+r)
	}

	total, err := models.UsersCollection.CountDocuments(context.Background(), bson.M{})

	return c.Status(fiber.StatusOK).JSON(rangeUser{
		Status:  "success",
		Offset:  offset,
		Limmit:  limmit,
		Total:   total,
		Count:   len(users),
		Next:    next,
		Message: "Sucessfuly found users",
		Users:   users,
	})
}
