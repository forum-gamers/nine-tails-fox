package database

import (
	"context"
	"log"
	"os"
	"time"

	h "github.com/forum-gamers/nine-tails-fox/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connection() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	h.PanicIfError(err)
	h.PanicIfError(client.Ping(ctx, readpref.Primary()))

	log.Println("Connected to the database")
	DB = client.Database("Post")
}
