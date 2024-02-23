package core

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database = mongo.Database

func InitDatabase() *Database {
	var database *Database
	config := Configuration()
	var URI string
	if config.DATABASE.USER != "" && config.DATABASE.PASSWORD != "" {
		URI = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.DATABASE.USER, config.DATABASE.PASSWORD, config.DATABASE.HOST, config.DATABASE.PORT)
	} else {
		URI = fmt.Sprintf("mongodb://%s:%s", config.DATABASE.HOST, config.DATABASE.PORT)
	}
	params := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(context.Background(), params)
	if err != nil {
		log.Fatalf("failed to connect to the database %s", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to ping the database %s", err)
	}
	database = client.Database(config.DATABASE.NAME)
	return database
}
