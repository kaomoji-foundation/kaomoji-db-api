package auth

// Common specific stuff for /auth operations goes here! UwU

import (
	stdMsg "GO-API-template/src/helpers/stdMessages"

	"github.com/gofiber/fiber/v2"
)

// Standard error response for user related authorization error responses //? shuld this be moved to stdMessages?
func stdAuthError(c *fiber.Ctx, status, message string, data interface{}) error {
	return c.Status(fiber.StatusUnauthorized).JSON(stdMsg.ErrorDefault(message, data))
}
