package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Setenv("JWT_SECRET_KEY", fmt.Sprintf("%x", NewRandomKey()))
	os.Setenv("ENCRYPT_SECRET_KEY", fmt.Sprintf("%x", NewRandomKey()))
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = ""
	}
	return value
}
