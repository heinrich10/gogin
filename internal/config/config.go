package config

import (
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	DB_HOST string
	HOST    string
	PORT    int
	API_KEY string
}

func LoadConfig() config {
	return config{
		DB_HOST: os.Getenv("DB_HOST"),
		HOST:    os.Getenv("HOST"),
		API_KEY: os.Getenv("API_KEY"),
		PORT: func() int {
			if s := os.Getenv("PORT"); s != "" {
				if v, err := strconv.Atoi(s); err == nil {
					return v
				}
			}
			return 0
		}(),
	}
}

var Config = LoadConfig()
