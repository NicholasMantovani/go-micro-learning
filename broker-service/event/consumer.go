package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
	"log"
	"net/http"
)

var loggerBaseUrl string

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) List(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err := ch.QueueBind(
			q.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message on exchange [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever // This blocks forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", loggerBaseUrl+"/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("cloud not close request body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
