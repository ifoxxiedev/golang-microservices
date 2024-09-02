package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	log.Printf("Starting server on port %s\n", webPort)
	rabbitConn, err := connect()

	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	defer rabbitConn.Close()

	app := Config{
		Rabbit: rabbitConn,
	}

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start http server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func connect() (*amqp.Connection, error) {
	// connect to rabbitmq
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until we have a connection
	for {
		c, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
		if err != nil {
			fmt.Println("RabbitMQ not yeat ready", backOff)
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(backOff), 2)) * time.Second
		log.Println("backing off", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
