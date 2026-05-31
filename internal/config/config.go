package config

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DB_HOST         string
	HOST            string
	PORT            int
	TRUSTED_PROXIES []string
	ALLOWED_ORIGINS []string
}

func LoadConfig() *Config {
	cfg := &Config{
		DB_HOST:         getenv[string]("DB_HOST", ""),
		HOST:            getenv[string]("HOST", ""),
		PORT:            getenv[int]("PORT", 3000),
		TRUSTED_PROXIES: getenv[[]string]("TRUSTED_PROXIES", []string{"127.0.0.1"}),
		ALLOWED_ORIGINS: getenv[[]string]("ALLOWED_ORIGINS", []string{"*"}),
	}
	return cfg
}

func getenv[T any](key string, def T) T {
	if s := os.Getenv(key); s != "" {
		switch any(def).(type) {
		case string:
			return any(s).(T)
		case int:
			if v, err := strconv.Atoi(s); err == nil {
				return any(v).(T)
			}
		case []string:
			parts := strings.Split(s, ",")
			var list []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					list = append(list, p)
				}
			}
			return any(list).(T)
		default:
			return def
		}
	}
	return def
}
