package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type TokenConfig struct {
	Secret string
	Issuer string
	Expire int
}

type Config struct {
	DatabaseUrl   string
	Host          string
	Port          int
	AllowedOrigin string

	token *TokenConfig
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env file not found, using system environment")
	}
	portNum, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}

	tokenExpire, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_MINUTES"))
	if err != nil {
		panic(err)
	}

	return &Config{
		DatabaseUrl:   os.Getenv("DATABASE_URL"),
		Host:          os.Getenv("HOST"),
		Port:          portNum,
		AllowedOrigin: os.Getenv("ALLOWED_ORIGIN"),
		token: &TokenConfig{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("JWT_ISSUER"),
			Expire: tokenExpire,
		},
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

func (c *Config) GetTokenConfig() *TokenConfig {
	return c.token
}
