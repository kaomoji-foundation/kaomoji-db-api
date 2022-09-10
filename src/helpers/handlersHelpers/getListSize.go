package handlershelpers

import (
	stdMsg "kaomojidb/src/helpers/stdMessages"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

/*
Gets the o (offset) and l (limmit) url query parameters
*/
func GetListSize(c *fiber.Ctx) (offset int, limmit int, err error) {
	offset, offsetErr := strconv.Atoi(c.Query("o", "0"))
	limmit, limmitErr := strconv.Atoi(c.Query("l", "10"))
	if (limmit - offset) > 100 {
		limmit = offset + 100
	}
	if offsetErr != nil || limmitErr != nil {
		offset = 0
		limmit = 10
		err = c.Status(fiber.StatusBadRequest).JSON(stdMsg.ErrorDefault(
			"Found an error while parsing your input, review your 'o' and 'l' query params", nil,
		))
	}
	return
}
