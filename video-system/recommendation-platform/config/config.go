package config

import (
	"os"
	"strings"
)

type Config struct {
	Port             string
	IndexBase        string
	JWTSecret        string
	ContentPlatforms map[string]string
	DBPath           string
}

func LoadConfig() *Config {

	return &Config{
		Port:             envString("RECOMMEND_PORT", "8082"),
		IndexBase:        envString("INDEX_BASE", "http://localhost:8083"),
		JWTSecret:        envString("JWT_SECRET", "dev-secret"),
		ContentPlatforms: envMap("CONTENT_PLATFORMS"),
		DBPath:           envString("DB_PATH", "./data/recommendation.db"),
	}

}

func envString(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envMap(key string) map[string]string {
	raw := os.Getenv(key)
	if raw == "" {
		return nil
	}
	out := map[string]string{}
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
