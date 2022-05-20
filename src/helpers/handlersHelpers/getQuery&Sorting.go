package handlershelpers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Query struct {
	Filter string `json:"query"`
	// search scopes [categories, description]
	Scopes []string `json:"scopes"`
}

type Sorting struct {
	ByPopularity           int8 `json:"byPopularity"`
	ByCategoriesAlphabetic int8 `json:"byCategoriesAlphabetic"`
}

//! Sorting not implemented, Filter (query) is
func GetQueryAndSorting(c *fiber.Ctx) (query Query, sorting Sorting, err error) {
	query.Filter = c.Query("query", "")
	query.Scopes = strings.Split(c.Query("scopes", "categories"), ",")

	//TODO: implement sorting

	/* 	if offsetErr != nil || limmitErr != nil {
		offset = 0
		limmit = 10
		err = c.Status(fiber.StatusBadRequest).JSON(stdMsg.ErrorDefault(
			"Found an error while parsing your input, review your 'o' and 'l' query params", nil,
		))
	} */
	return
}
