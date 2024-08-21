package config

import "github.com/joho/godotenv"

func ConfigENV() error {
	err := godotenv.Load()

	return err
}
