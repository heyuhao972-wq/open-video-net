package config

import "os"

type RuntimeConfig struct {
	HTTPPort string
	ChunkDir string
}

func LoadRuntimeConfig() *RuntimeConfig {
	return &RuntimeConfig{
		HTTPPort: envString("P2P_HTTP_PORT", "8090"),
		ChunkDir: envString("CHUNK_DIR", "./data/chunks"),
	}
}

func envString(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
