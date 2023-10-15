package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"main/services"
	"os"
)

type JSON gin.H

var (
	Postgres        *pg.DB
	MongoCollection *mongo.Collection
	MongoClient     *mongo.Client
)

func init() {
	// Parse env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln(err)
	}

	// Connect to postgres
	Postgres = services.ConnectToPostgres(os.Getenv("PG_DATABASE"))

	// connect to Mongo
	MongoCollection, MongoClient = services.ConnectToMongo()

	// file name and line in logs
	log.SetFlags(log.LstdFlags | log.Llongfile)

	// Set gin to release
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("Init")
}

func main() {
	defer Postgres.Close()
	defer MongoClient.Disconnect(context.Background())

	engine := gin.Default()
	engine.POST("/handler", Handler)
	err := engine.Run(os.Getenv("LISTEN_ADDR"))

	panic(err)
}
