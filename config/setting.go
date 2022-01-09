package config

import (
	"log"
	"os"
	"strconv"
)

type Setting struct {
	ServiceName        string
	GRPCServerEndpoint string
	GRPCServerPort     string
	HTTPServerPort     string

	ShouldProfile bool
	ShouldTrace   bool

	Env                       string
	GracefulShutdownTimeoutMs int
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if defaultValue == "" {
		log.Fatalf("a required environment variable missed: %s", key)
	}
	return defaultValue
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Panic(err)
	}
	return i
}

func mustAtob(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Panic(err)
	}
	return b
}

func NewSetting() Setting {
	return Setting{
		ServiceName:        "auth",
		GRPCServerEndpoint: getEnv("GRPC_SERVER_ENDPOINT", "localhost:18081"),
		GRPCServerPort:     getEnv("GRPC_SERVER_PORT", "18081"),
		HTTPServerPort:     getEnv("HTTP_SERVER_PORT", "18082"),

		ShouldProfile: mustAtob(getEnv("SHOULD_PROFILE", "false")),
		ShouldTrace:   mustAtob(getEnv("SHOULD_TRACE", "false")),

		Env:                       getEnv("ENV", "development"),
		GracefulShutdownTimeoutMs: mustAtoi(getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_MS", "5000")),
	}
}
