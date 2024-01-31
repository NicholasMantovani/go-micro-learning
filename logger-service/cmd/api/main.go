package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger/data"
	"net/http"
	"os"
	"time"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongodb
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	//go app.serve()
	log.Println("Starting service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

//func (app *Config) serve() {
//	srv := &http.Server{
//		Addr:    fmt.Sprintf(":%s", webPort),
//		Handler: app.routes(),
//	}
//
//	err := srv.ListenAndServe()
//	if err != nil {
//		log.Panic(err)
//	}
//}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options
	mongoUrl := os.Getenv("MONGODB_CONN")
	mongoUser := os.Getenv("MONGODB_USER")
	mongoPass := os.Getenv("MONGODB_PASSWORD")

	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: mongoUser,
		Password: mongoPass,
	})

	// connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo")

	return c, nil
}
