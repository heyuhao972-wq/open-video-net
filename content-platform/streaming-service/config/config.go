package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	PlatformBase   string
	P2PBase        string
	PlatformMap    map[string]string
	P2PMap         map[string]string
	ChunkTimeoutMs int
	ChunkRetry     int
	MaxParallel    int
	CacheSize      int
}

func LoadConfig() *Config {

	return &Config{
		Port:           envString("STREAM_PORT", "8081"),
		PlatformBase:   envString("PLATFORM_BASE", "http://localhost:8080"),
		P2PBase:        envString("P2P_BASE", "http://localhost:8090"),
		PlatformMap:    envMap("PLATFORM_MAP"),
		P2PMap:         envMap("P2P_MAP"),
		ChunkTimeoutMs: envInt("CHUNK_TIMEOUT_MS", 4000),
		ChunkRetry:     envInt("CHUNK_RETRY", 2),
		MaxParallel:    envInt("CHUNK_MAX_PARALLEL", 6),
		CacheSize:      envInt("CHUNK_CACHE_SIZE", 128),
	}

}

func envString(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
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
