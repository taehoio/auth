package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Setting struct {
	ServiceName        string
	GRPCServerEndpoint string
	GRPCServerPort     string
	HTTPServerPort     string

	Env                       string
	GracefulShutdownTimeoutMs int

	ShouldProfile bool
	ShouldTrace   bool

	JWTHMACSecret           string
	JWTIssuer               string
	JWTAudience             string
	AccessTokenExpiresInMs  int
	RefreshTokenExpiresInMs int
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

		Env:                       getEnv("ENV", "development"),
		GracefulShutdownTimeoutMs: mustAtoi(getEnv("GRACEFUL_SHUTDOWN_TIMEOUT_MS", "5000")),

		ShouldProfile: mustAtob(getEnv("SHOULD_PROFILE", "false")),
		ShouldTrace:   mustAtob(getEnv("SHOULD_TRACE", "false")),

		JWTHMACSecret:           getEnv("JWT_HMAC_SECRET", "PLEASE_SET_THIS_ENV_VAR"),
		JWTIssuer:               getEnv("JWT_ISSUER", "taeho.io"),
		JWTAudience:             getEnv("JWT_AUDIENCE", "taeho.io"),
		AccessTokenExpiresInMs:  mustAtoi(getEnv("ACCESS_TOKEN_EXPIRES_IN_MS", fmt.Sprintf("%d", 1000*60*15))),
		RefreshTokenExpiresInMs: mustAtoi(getEnv("REFRESH_TOKEN_EXPIRES_IN_MS", fmt.Sprintf("%d", 1000*60*60*24*365))),
	}
}
