package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connectToRabbitMq()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create a consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Panicln(err)
	}

	event.SetLoggerBaseUrl(os.Getenv("LOGGER_BASE_URL"))

	// watch the queue and consume events from topic
	err = consumer.List([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

}

func connectToRabbitMq() (*amqp.Connection, error) {
	// backoff rutine
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(os.Getenv("AMQP_CONNECTION"))
		if err != nil {
			log.Println("RabbitMq not yet ready..")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("backing off for %s...\n", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
