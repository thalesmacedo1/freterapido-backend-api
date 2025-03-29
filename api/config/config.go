package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Postgres Configuration
	PostgresURI      string
	PostgresUser     string
	PostgresPassword string

	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

var Settings Config

func LoadConfig(envFile string) error {
	if err := godotenv.Load(envFile); err != nil {
		return err
	}

	// Carregar as configurações do Postgres
	Settings.PostgresURI = strings.TrimSpace(os.Getenv("POSTGRES_URI"))
	Settings.PostgresUser = strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	Settings.PostgresPassword = strings.TrimSpace(os.Getenv("POSTGRES_PASSWORD"))

	// Validar configurações do Postgres
	if Settings.PostgresURI == "" || Settings.PostgresUser == "" || Settings.PostgresPassword == "" {
		log.Fatal("Missing required Postgres environment variables")
	}

	// Carregar as configurações do Redis
	Settings.RedisHost = strings.TrimSpace(os.Getenv("REDIS_HOST"))
	Settings.RedisPort = strings.TrimSpace(os.Getenv("REDIS_PORT"))
	Settings.RedisPassword = strings.TrimSpace(os.Getenv("REDIS_PASSWORD"))
	redisDBStr := strings.TrimSpace(os.Getenv("REDIS_DB"))

	// Validar configurações do Redis
	if Settings.RedisHost == "" || Settings.RedisPort == "" {
		log.Fatal("Missing required Redis environment variables (REDIS_HOST, REDIS_PORT)")
	}

	// Converter RedisDB para int
	if redisDBStr == "" {
		Settings.RedisDB = 0
	} else {
		db, err := strconv.Atoi(redisDBStr)
		if err != nil {
			log.Fatalf("Invalid Redis DB value: %v", err)
		}
		Settings.RedisDB = db
	}

	return nil
}
