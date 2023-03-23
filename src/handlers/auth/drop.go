package auth

import (
	"context"
	stdMsg "kaomojidb/src/helpers/stdMessages"
	"kaomojidb/src/models"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// Drop get jwt (get function so it "gets dropped")
// @Summary      drops (logout) bearer token
// @Description  Blocks a given token, if none is provided, the one used to acess this route willbe blocked.
// @Accept       json
// @Produce      json
// @Success      200  {object}  interface{}
// @Failure      403  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /auth/drop [get]
func Drop(c *fiber.Ctx) error {

	type reqBodyforToken struct {
		Token string `json:"token"`
	}

	var reqBody reqBodyforToken

	var userData models.User

	if len(c.Body()) != 0 { //empty body
		if err := c.BodyParser(&reqBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Error on drop request",
				"data":    err,
				"expected": fiber.Map{
					"token": "JWT bearer token",
				},
			})
		}
		reqBody.Token = ""
	}

	// Token of the user
	token := c.Locals("user").(*jwt.Token)
	tokenStr := token.Raw
	userID := token.Claims.(jwt.MapClaims)["uid"].(string)

	err := userData.Fill(userID, true, false, false)
	if err != nil {
		return stdAuthError(c, "error", "User not found", err)
	}
	// END get userdata from db

	// data sanity check
	if userData.Tokens == nil {
		userData.Tokens = map[string]bool{}
	}
	if userData.BlockedTokens == nil {
		userData.BlockedTokens = map[string]bool{}
	}

	if reqBody.Token != "" {
		tokenStr = reqBody.Token
	}

	// check if the token attempted to drop is on the Tokens map
	if _, ok := userData.Tokens[tokenStr]; !ok {
		return stdAuthError(c,
			"error",
			"Specified token does not exist or has already been dropped",
			reqBody,
		)
	}

	// block the token
	userData.BlockedTokens[tokenStr] = true
	delete(userData.Tokens, tokenStr)

	// prune expired tokens
	go userData.PruneTokens()

	filter := bson.M{"_id": userData.ID}
	update := bson.M{"$set": userData}
	_, err = models.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(stdMsg.ErrorDefault("An error ocurred while procesing the request", err))
	}

	return c.JSON(fiber.Map{"status": "success", "message": "token dropped sucessfully", "token": tokenStr})
}
