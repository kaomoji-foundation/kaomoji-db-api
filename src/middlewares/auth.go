package middlewares

import (
	cfg "GO-API-template/src/config"
	stdMsg "GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Protected protect routes
func Auth() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(cfg.Config.JWT.Secret),
		ErrorHandler:   jwtError,
		SuccessHandler: TokenCheck,
	})
}

// Optional Authentication, sets var to true or false
func OptInAuth() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(cfg.Config.JWT.Secret),
		ErrorHandler:   jwtInvalid,
		SuccessHandler: TokenCheck,
	})
}

// ckecks a token to be valid against it's user
func TokenCheck(c *fiber.Ctx) error {

	// Get Token of the reader's user
	tok := c.Locals("user")
	token := tok.(*jwt.Token)
	uid := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["uid"])

	var user models.User
	userID, err := primitive.ObjectIDFromHex(uid)
	user.ID = userID
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(stdMsg.ErrorDefault("An error ocurred while Authenticating the request", nil))
	}

	err = user.LoadTokens()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(stdMsg.ErrorDefault("An error ocurred while Authenticating the request", nil))
	}

	allowed := !user.BlockedTokens[token.Raw]
	if !allowed {
		return jwtError(c, nil)
	}

	c.Locals("Authenticated", true)
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	const ErrJWTMissingOrMalformed = "Missing or malformed JWT"
	if err != nil && err.Error() != ErrJWTMissingOrMalformed {
		return c.Status(fiber.StatusUnauthorized).
			JSON(stdMsg.ErrorDefault(ErrJWTMissingOrMalformed, nil))
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(stdMsg.ErrorDefault("Invalid or expired JWT", nil))
}

// ckecks a token to be valid against it's user, if not inform of it and continue
func jwtInvalid(c *fiber.Ctx, err error) error {
	c.Locals("Authenticated", false)
	const ErrJWTMissingOrMalformed = "Missing or malformed JWT"
	if err != nil && err.Error() == ErrJWTMissingOrMalformed {
		return c.Next()
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(stdMsg.ErrorDefault("Invalid or expired JWT", nil))
}
