package config

import (
	"os"
)

type mongoConfig struct {
	URI  string
	Pass string
	User string
	DBs  []string
}

var mongoDB mongoConfig
var loadedMongo bool

func MongoDB(reload ...bool) *mongoConfig {
	if loadedMongo && !reload[0] {
		return &mongoDB
	}

	mongoDB = mongoConfig{
		URI:  os.Getenv("DB_URI"),
		Pass: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBs: []string{
			os.Getenv("DB_1_NAME"),
		},
	}

	// Defaults:
	//! BE CAREFULL can be dangerous on certain enviroments if not controlled propertly

	if mongoDB.User == "" {
		mongoDB.User = "root" //default to a local development db credentials
	}

	if mongoDB.Pass == "" {
		mongoDB.Pass = "example" //default to a local development db credentials
	}

	if mongoDB.URI == "" {
		mongoDB.URI = "mongodb://localhost:27017" //default to a local development db
	}

	if mongoDB.DBs[0] == "" {
		mongoDB.DBs[0] = "devel-kaomoji-DB"
	}
	loadedMongo = true
	return &mongoDB
}
