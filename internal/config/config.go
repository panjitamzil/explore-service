package config

import (
	"os"
)

type Config struct {
	GRPCPort  string
	MySQLDSN  string
	RedisAddr string
}

func FromEnv() Config {
	cfg := Config{
		GRPCPort:  getEnv("GRPC_PORT", "50051"),
		MySQLDSN:  getEnv("MYSQL_DSN", "root:password@tcp(mysql:3306)/explore?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"),
		RedisAddr: getEnv("REDIS_ADDR", "redis:6379"),
	}
	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
