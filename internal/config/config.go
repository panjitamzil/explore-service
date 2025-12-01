package config

import (
	"os"
)

type Config struct {
	GRPCPort string
}

func FromEnv() Config {
	cfg := Config{
		GRPCPort: getEnv("GRPC_PORT", "50051"),
	}
	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
