package kaomojis

import (
	hh "GO-API-template/src/helpers/handlersHelpers"
	stdMsg "GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"GO-API-template/src/utils/filtering"
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type rangeKaomoji struct {
	Status   string      `json:"status"`
	Message  string      `json:"message"`
	Total    int64       `json:"total"`
	Count    int         `json:"count"`
	Offset   int         `json:"offset"`
	Limmit   int         `json:"limmit"`
	Next     string      `json:"next"`
	Kaomojis interface{} `json:"kaomojis"`
}

// GetKaomojis get the kaomoji list
// @Summary      Retrieve kaomoji list
// @Description  Retrieve the kaomoji list
// @security     BearerAuth
// @Accept       json
// @Produce      json
// @Router       /kaomojis [get]
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      404  {object}  interface{}
// @Failure      500  {object}  interface{}
func GetKaomojis(c *fiber.Ctx) error {
	offset, limmit, err := hh.GetListSize(c)
	if err != nil {
		return err
	}
	// sorting is not implemented in GetQueryAndSorting
	query, _, err := hh.GetQueryAndSorting(c)
	if err != nil {
		return err
	}

	// to match models.kaomojiMinimal
	projection := bson.M{"_id": 1, "string": 1, "description": 1, "categories": 1}

	var opts = options.Find().SetProjection(projection)
	// If searching is required, it's done here, so we need all the data from the db to search from
	if query.Filter == "" {
		opts = opts.SetSkip(int64(offset)).SetLimit(int64(limmit))
	}

	cursor, err := models.KaomojisCollection.Find(context.Background(), bson.D{}, opts)
	if err != nil && err != mongo.ErrNoDocuments {
		c.Status(fiber.StatusInternalServerError).JSON(
			stdMsg.ErrorDefault("An error ocurred while retrieving the kaomojis data", nil),
		)
	}
	defer cursor.Close(context.Background())

	var kaomojis []models.KaomojiMinimal
	type searchRessults struct {
		// Represents the amount of coincidences within the kaomoji usefull for later sorting
		// The lower index of data the better the result is
		// data[0][n] is a better result than data[1][n]
		data [][]models.KaomojiMinimal
	}

	// total of unfiltered documents //// 0 if filtering
	total := int64(0)
	kaoRange := limmit - offset

	query.Filter = strings.TrimSpace(query.Filter)

	if query.Filter == "" {
		err = cursor.All(context.Background(), &kaomojis)
	} else {
		r := searchRessults{
			data: make([][]models.KaomojiMinimal, 3), // [0:[], 1:[], 2:[]]
		}
		// processes every kaomoji, but only adds to kaomjis if still in range
		for i := 0; cursor.Next(context.Background()); i++ {
			var k models.KaomojiMinimal
			if err = cursor.Decode(&k); err != nil {
				c.Status(fiber.StatusInternalServerError).JSON(
					stdMsg.ErrorDefault("An error ocurred while retrieving the kaomojis data", nil),
				)
			}

			// data optimization for search
			var categories = make([]string, 0, len(k.Categories))
			for _, c := range k.Categories {
				if subCat := strings.Split(c, " "); len(subCat) > 1 {
					for _, sc := range subCat {
						categories = append(categories, sc)
					}
				}
				categories = append(categories, c)
			}
			filters := strings.Split(query.Filter, " ")

			// filter valid Categories for the given filters (taken from spliting query.filter)
			filteredCategories, distances := filtering.SearchStrings(categories, filters, 2)

			// ckecks k is a valid result and not out of range for this request
			if len(filteredCategories) > 0 && len(kaomojis) <= kaoRange {
				totalDist := 0
				for _, distance := range distances {
					totalDist += distance
				}
				averageDistance := totalDist / len(distances)

				r.data[averageDistance] = append(r.data[averageDistance], k)
			}
			if len(filteredCategories) > 0 {
				total++
			}
		}

		// append the ranked results from best to worst
		for _, d := range r.data {
			kaomojis = append(kaomojis, d...)
		}
	}

	var next string
	if kaoRange <= len(kaomojis) {
		next = c.BaseURL() + string(c.Request().URI().Path()) + fmt.Sprintf("/?o=%v&l=%v", offset+kaoRange, limmit+kaoRange)
	}

	// total of unfiltered documents, 0 if filtering
	if total == 0 {
		total, err = models.KaomojisCollection.CountDocuments(context.Background(), bson.M{})
	}

	//TODO: customize Message for when filtering
	return c.Status(fiber.StatusOK).JSON(rangeKaomoji{
		Status:   "success",
		Offset:   offset,
		Limmit:   limmit,
		Total:    total,
		Count:    len(kaomojis),
		Next:     next,
		Message:  "Sucessfuly found kaomojis",
		Kaomojis: kaomojis,
	})
}
