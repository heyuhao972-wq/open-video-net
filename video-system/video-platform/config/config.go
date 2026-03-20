package config

import (
	"os"
	"strings"
)

type Config struct {
	Port       string
	JWTSecret  string
	IndexBase  string
	PlatformID string
	AcceptTags []string
	DBPath     string
}

func LoadConfig() *Config {

	return &Config{
		Port:       envString("PLATFORM_PORT", "8080"),
		JWTSecret:  envString("JWT_SECRET", "dev-secret"),
		IndexBase:  envString("INDEX_BASE", ""),
		PlatformID: envString("PLATFORM_ID", "platformA"),
		AcceptTags: envList("ACCEPT_TAGS"),
		DBPath:     envString("DB_PATH", "./data/platform.db"),
	}

}

func envString(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envList(key string) []string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}
