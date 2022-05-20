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
		data         *[]models.KaomojiMinimal
		resultRating []int // represents the amount of coincidences within the kaomoji usefull for later sorting
	}

	// total of unfiltered documents, 0 if filtering
	total := int64(0)
	kaoRange := limmit - offset

	if query.Filter == "" {
		err = cursor.All(context.Background(), &kaomojis)
	} else {
		r := searchRessults{
			data:         &kaomojis,
			resultRating: make([]int, 0, len(kaomojis)),
		}
		// processes every kaomoji, but only adds to kaomjis if still in range
		for cursor.Next(context.Background()) {
			var k models.KaomojiMinimal
			if err = cursor.Decode(&k); err != nil {
				c.Status(fiber.StatusInternalServerError).JSON(
					stdMsg.ErrorDefault("An error ocurred while retrieving the kaomojis data", nil),
				)
			}

			// data optimization for search
			var cat = make([]string, 0, len(k.Categories))
			for i := 0; i < len(k.Categories); i++ {
				c := k.Categories[i]
				if subCat := strings.Split(c, " "); len(subCat) > 1 {
					for i := 0; i < len(subCat); i++ {
						cat = append(cat, subCat[i])
					}
				}
				cat = append(cat, c)
			}

			// filtered Categories
			fcat := filtering.SearchString(cat, query.Filter)

			// ckecks k is a valid result and not out of range for this request
			if len(fcat) > 0 && len(kaomojis) <= kaoRange {
				r.resultRating = append(r.resultRating, len(fcat))
				//kaomojis[len(kaomojis)] = k
				*r.data = append(*r.data, k)
			}
			if len(fcat) > 0 {
				total++
			}
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
