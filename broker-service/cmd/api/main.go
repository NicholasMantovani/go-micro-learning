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

const webPort = "8080"

type Config struct {
	AuthBaseUrl   string
	LoggerDNSName string
	MailerBaseUrl string
	Rabbit        *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connectToRabbitMq()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()

	authBaseUrl := os.Getenv("AUTH_BASE_URL")
	loggerDnsName := os.Getenv("LOGGER_DNS_NAME")
	mailerBaseUrl := os.Getenv("MAILER_BASE_URL")

	app := Config{
		AuthBaseUrl:   authBaseUrl,
		LoggerDNSName: loggerDnsName,
		MailerBaseUrl: mailerBaseUrl,
		Rabbit:        rabbitConn,
	}
	log.Printf("Starting broken server on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the serer
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
