package config

import "os"

type bearerJWTConfig struct {
	Secret string
}

var confJWT bearerJWTConfig

func JWT(reload ...bool) *bearerJWTConfig {
	if (bearerJWTConfig{}) != confJWT && !reload[0] {
		return &confJWT
	}
	confJWT = bearerJWTConfig{
		Secret: os.Getenv("JWT_SECRET"),
	}
	return &confJWT
}
