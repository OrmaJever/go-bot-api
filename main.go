package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type JSON gin.H

var Postgres *pg.DB
var MongoCollection *mongo.Collection

func init() {
	// Parse env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln(err)
	}

	// Connect to postgres
	Postgres = pg.Connect(&pg.Options{
		Addr:     os.Getenv("PG_ADDR"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Database: os.Getenv("PG_DATABASE"),
	})

	// connect to Mongo
	MongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION")))

	if err != nil {
		log.Fatalln(err)
	}

	MongoCollection = MongoClient.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	// Set log stream to file
	logFile, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set gin to release
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("Init")
}

func main() {
	engine := gin.Default()

	engine.POST("/handler", Handler)
	engine.Run(os.Getenv("LISTEN_ADDR"))
}
