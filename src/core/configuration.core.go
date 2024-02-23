package core

import (
	"log"
	"os"
	"strconv"

	dotEnv "github.com/joho/godotenv"
)

type MailerConfig struct {
	FROM string
}

type DatabaseConfig struct {
	HOST     string
	PORT     string
	USER     string
	PASSWORD string
	NAME     string
}

type Config struct {
	PORT       int
	DATABASE   DatabaseConfig
	STAGE      string
	AUDIENCE   string
	JWT_SECRET string
	JWT_EXPIRY int
	MAILER     MailerConfig
}

var instance *Config

func Configuration() *Config {
	if instance == nil {
		err := dotEnv.Load()
		if err != nil {
			log.Fatalf("failed to load .env file %s", err)
		}
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			port = 8080
		}

		db_port := os.Getenv("MONGODB_PORT")
		db_host := os.Getenv("MONGODB_HOST")
		db_user := os.Getenv("MONGODB_USER")
		db_password := os.Getenv("MONGODB_PASSWORD")
		db_name := os.Getenv("MONGODB_NAME")
		audience := os.Getenv("AUDIENCE")
		stage := os.Getenv("ENV")
		jwt_secret := os.Getenv("JWT_SECRET")
		mailer_from := os.Getenv("MAILER_FROM")

		if stage == "" {
			stage = "development"
		}
		if jwt_secret == "" {
			log.Fatalf("jwt secret not set")
		}
		jwt_expiry, err := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
		if err != nil {
			jwt_expiry = 24
		}

		instance = &Config{
			PORT: port,
			DATABASE: DatabaseConfig{
				PORT:     db_port,
				HOST:     db_host,
				USER:     db_user,
				PASSWORD: db_password,
				NAME:     db_name,
			},
			STAGE:      stage,
			AUDIENCE:   audience,
			JWT_SECRET: jwt_secret,
			JWT_EXPIRY: jwt_expiry,
			MAILER: MailerConfig{
				FROM: mailer_from,
			},
		}
	}
	return instance
}
