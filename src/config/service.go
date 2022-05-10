package config

import (
	"os"
)

type serviceConfig struct {
	Port string
}

var service serviceConfig

func Service(reload ...bool) *serviceConfig {
	if (serviceConfig{}) != service && !reload[0] {
		return &service
	}

	service.Port = os.Getenv("PORT")
	if service.Port == "" {
		service.Port = "3000"
	}
	return &service
}
