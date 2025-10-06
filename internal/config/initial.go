package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func getEnvInt(val string) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("failed to convert to integer: %w", err))
	}

	return res
}

func Load(appEnv string) *Config {
	err := godotenv.Load()

	if appEnv == "test" {
		err = godotenv.Load(filepath.Join("..", ".env.test"))
	}

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &Config{
		App: App{
			Name: os.Getenv("APP_NAME"),
		},
		Server: Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
			Prefork: os.Getenv("SERVER_PREFORK"),
		},
		Database: Database{
			Host: os.Getenv("DATABASE_HOST"),
			Port: os.Getenv("DATABASE_PORT"),
			Name: os.Getenv("DATABASE_NAME"),
			User: os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Timezone: os.Getenv("DATABASE_TIMEZONE"),
			Pool: Pool{
				Idle: getEnvInt(os.Getenv("POOL_IDLE")),
				Max: getEnvInt(os.Getenv("POOL_MAX")),
				Lifetime: getEnvInt(os.Getenv("POOL_LIFETIME")),
			},
		},
		Logrus: Logrus{
			Level: int32(getEnvInt(os.Getenv("LOG_LEVEL"))),
		},
	}
}