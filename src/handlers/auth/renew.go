package auth

import (
	cfg "GO-API-template/src/config"
	stdMsg "GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// Login get user and password
// @Summary      Renew the authentication bearer token, blocks the token acessed to use this
// @Description  Get your new user's token to acess users only protected routes
// @Accept       json
// @Produce      json
// @Success      200  {object}  interface{}
// @Failure      403  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /auth/renew [get]
func Renew(c *fiber.Ctx) error {
	var userData models.User
	// Token of the user
	token := c.Locals("user").(*jwt.Token)
	userID := token.Claims.(jwt.MapClaims)["uid"].(string)

	err := userData.Fill(userID, true, false, false)
	if err != nil {
		return stdAuthError(c, "error", "User not found", err)
	}
	// END get userdata from db

	// create the new claims
	claims := jwt.MapClaims{
		"username": userData.Username,
		"email":    userData.Email,
		"uid":      userData.ID.Hex(),
		"role":     userData.Role,
		"rid":      userData.RoleID.Hex(),
		"exp":      time.Now().Add(time.Hour * 24 * 5).Unix(),
	}

	// create token
	tokenRenewed := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	t, err := tokenRenewed.SignedString([]byte(cfg.Config.JWT.Secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if userData.Tokens == nil {
		userData.Tokens = map[string]bool{}
	}
	if userData.BlockedTokens == nil {
		userData.BlockedTokens = map[string]bool{}
	}
	userData.BlockedTokens[token.Raw] = true
	delete(userData.Tokens, token.Raw)
	userData.Tokens[t] = true

	filter := bson.M{"_id": userData.ID}
	update := bson.M{"$set": userData}
	_, err = models.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(stdMsg.ErrorDefault("An error ocurred while procesing the request", err))
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Login sucessful", "token": t})
}
