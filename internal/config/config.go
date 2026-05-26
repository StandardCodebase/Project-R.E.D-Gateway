package config

import "os"

type Config struct {
	NodeName string
	DataDir  string
	Port     string
}

func Load() Config {
	return Config{
		NodeName: envOr("RED_NODE_NAME", "Alpha-Centauri-01"),
		DataDir:  envOr("RED_DATA_DIR", "./data"),
		Port:     envOr("RED_PORT", "8080"),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
