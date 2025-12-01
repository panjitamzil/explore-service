package config

import (
	"os"
)

type Config struct {
	GRPCPort string

	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDB       string

	RedisHost string
	RedisPort string
}

func FromEnv() Config {
	cfg := Config{
		GRPCPort:      getEnv("GRPC_PORT", "50051"),
		MySQLHost:     getEnv("MYSQL_HOST", "mysql"),
		MySQLPort:     getEnv("MYSQL_PORT", "3306"),
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", "password"),
		MySQLDB:       getEnv("MYSQL_DB", "explore"),
		RedisHost:     getEnv("REDIS_HOST", "redis"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
	}
	return cfg
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

func (c Config) MySQLDSN() string {
	return c.MySQLUser + ":" + c.MySQLPassword + "@tcp(" + c.MySQLHost + ":" + c.MySQLPort + ")/" + c.MySQLDB + "?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci"
}

func (c Config) RedisAddr() string {
	return c.RedisHost + ":" + c.RedisPort
}
