package main

import (
	"context"
	"fmt"
	"log"
	"log-service/cmd/data"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort               = "80"
	rpcPort               = "5001"
	maxConnectionAttempts = 10
	gRpcPort              = "50001"
)

var mongoURL string
var client *mongo.Client
var counts int64

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoURL = os.Getenv("MONGO_URL")
	mongoClient, err := connectToDB()

	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	models := data.New(client)
	config := Config{Models: models}

	// Register RPC Server
	err = rpc.Register(new(RPCServer))
	go config.rpcListen()

	// Registert gRPC Server
	go config.gRPCListen()

	// Register HTTP Server
	config.httoListen()

}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC Server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}

	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) httoListen() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	fmt.Printf("Starting service on port %s\n", webPort)
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openMongoClient() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func connectToDB() (*mongo.Client, error) {
	for {
		client, err := openMongoClient()
		if err != nil {
			log.Println("MongoDB is not ready yet...", err)
			counts += 1
		} else {
			log.Println("MongoDB is ready!")
			return client, nil
		}

		if counts >= maxConnectionAttempts {
			log.Printf("MongoDB is not ready after 10 attempts. Exiting... %v\n", err)
			return nil, err
		}

		log.Println("Backing off for 5 seconds before trying again...")
		time.Sleep(1 * time.Second)
		continue
	}
}
