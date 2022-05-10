package config

import "os"

type openTelemetryConfig struct {
	LightStepKey string
}

var openTel openTelemetryConfig

func OpenTel(reload ...bool) *openTelemetryConfig {
	if openTel != (openTelemetryConfig{}) && !reload[0] {
		return &openTel
	}
	openTel = openTelemetryConfig{
		LightStepKey: os.Getenv("OTEL-LIGHT-STEP-KEY"),
	}
	return &openTel
}
