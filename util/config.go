package util

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver      string
	DBSource      string
	ServerAddress string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	var config Config
	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DBSource = os.Getenv("DB_SOURCE")
	config.ServerAddress = os.Getenv("SERVER_ADDRESS")

	return config
}
