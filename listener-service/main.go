package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		exitApp(err)
	}

	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
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

func exitApp(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
