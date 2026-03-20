package config

import "os"

type Config struct {
	Port      string
	IndexBase string
	JWTSecret string
}

func LoadConfig() *Config {

	return &Config{
		Port:      envString("RECOMMEND_PORT", "8082"),
		IndexBase: envString("INDEX_BASE", "http://localhost:8083"),
		JWTSecret: envString("JWT_SECRET", "dev-secret"),
	}

}

func envString(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
