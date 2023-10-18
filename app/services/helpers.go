package services

import (
	"context"
	"github.com/go-pg/pg/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"main/models"
	"main/telegram"
	"os"
	"time"
)

type CallbackT func(data *telegram.Data, tgApi *Telegram, bot *models.Bot)

func ConnectToMongo() (*mongo.Collection, *mongo.Client) {
	MongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION")))

	if err != nil {
		log.Fatalln(err)
	}

	var result bson.M
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	if err = MongoClient.Database(os.Getenv("MONGO_DB")).RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		log.Fatalln(err)
	}

	return MongoClient.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION")), MongoClient
}

func ConnectToPostgres(database string) *pg.DB {
	postgres := pg.Connect(&pg.Options{
		Addr:     os.Getenv("PG_ADDR"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Database: database,
	})

	_, err := postgres.Exec("SELECT 1")

	if err != nil {
		log.Fatalln(err)
	}

	if os.Getenv("GIN_MODE") == "debug" {
		postgres.AddQueryHook(PostgresLogger{})
	}

	return postgres

}
