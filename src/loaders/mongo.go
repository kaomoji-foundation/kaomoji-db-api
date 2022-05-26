package loaders

import (
	"context"

	"Kaomoji-DB/src/models"
	"Kaomoji-DB/src/services"
)

func LoadMongo() *context.CancelFunc {
	cancelCtx := services.Mongo.Init()

	// init all the collections
	models.User{}.CreateSingletonDBAndCollection()
	models.Role{}.CreateSingletonDBAndCollection()
	models.Kaomoji{}.CreateSingletonDBAndCollection()
	models.Issue{}.CreateSingletonDBAndCollection()
	return cancelCtx
}
