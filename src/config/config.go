package config

import (
	"github.com/joho/godotenv"

	"log"
)

type Configuration struct {
	Service *serviceConfig
	Mongo   *mongoConfig
	JWT     *bearerJWTConfig
	OpenTel *openTelemetryConfig
}

var Config Configuration

// loads all the config onto Configuration struct
func Load(forceReload ...bool) *Configuration {
	if forceReload[0] {
		LoadEnvFile()
	}
	Config.Service = Service(forceReload[0])
	Config.Mongo = MongoDB(forceReload[0])
	Config.JWT = JWT(forceReload[0])
	Config.OpenTel = OpenTel(forceReload[0])

	return &Config
}

// dinamically get the variables, so
func LoadEnvFile() bool {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Could not load ENV vars from file")
		return false
	}
	return true
}
