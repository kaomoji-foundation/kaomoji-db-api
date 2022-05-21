package config

import (
	"os"
	"strconv"
)

type serviceConfig struct {
	Port string
	// Time in seconds
	LimitorTimeFrame int
	LimitorLimit     int
}

var service serviceConfig

func Service(reload ...bool) *serviceConfig {
	if (serviceConfig{}) != service && !reload[0] {
		return &service
	}

	service.Port = os.Getenv("PORT")
	service.LimitorTimeFrame, _ = strconv.Atoi(os.Getenv("LimitorTimeFrame"))
	service.LimitorLimit, _ = strconv.Atoi(os.Getenv("LimitorLimit"))
	if service.Port == "" {
		service.Port = "3000"
	}
	if service.LimitorTimeFrame == 0 {
		service.LimitorTimeFrame = 10
	}
	if service.LimitorLimit == 0 {
		service.LimitorLimit = 20
	}
	return &service
}
