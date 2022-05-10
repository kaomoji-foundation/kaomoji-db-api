package auth

import (
	cfg "GO-API-template/src/config"
	stdMsg "GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"GO-API-template/src/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// Login get user and password
// @Summary      login to get the authentication bearer token
// @Description  Get your user's token to acess users only protected routes
// @Accept       json
// @Produce      json
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      400  {object}  interface{}
// @Router       /auth/login [get]
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	var input LoginInput
	var userData models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Error on login request",
			"data":    err,
			"expected": fiber.Map{
				"identity": "username or password",
				"password": "user password",
			},
		})
	}

	err := userData.Fill(input.Identity, false, true, true)
	if err != nil {
		return stdAuthError(c, "error", "User not found", err)
	}
	// END get userdata from db

	// validate provided credenntials for user
	if !utils.CheckHash(input.Password, userData.Password) {
		return stdAuthError(c, "error", "Invalid password", nil)
	}
	// END validate provided credenntials for user

	// create the claims
	claims := jwt.MapClaims{
		"username": userData.Username,
		"email":    userData.Email,
		"uid":      userData.ID.Hex(),
		"role":     userData.Role,
		"rid":      userData.RoleID.Hex(),
		"exp":      time.Now().Add(time.Hour * 24 * 5).Unix(),
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte(cfg.Config.JWT.Secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if userData.Tokens == nil {
		userData.Tokens = map[string]bool{}
	}
	userData.Tokens[t] = true

	filter := bson.M{"_id": userData.ID}
	update := bson.M{"$set": userData}
	_, err = models.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(stdMsg.ErrorDefault("An error ocurred while procesing the request", err))
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Login sucessful", "token": t})
}
