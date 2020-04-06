package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type consumeHandler struct {
	consumer *consumer
	upgrader *websocket.Upgrader
}

func newConsumeHandler(consumer *consumer, upgrader *websocket.Upgrader) *consumeHandler {
	return &consumeHandler{consumer: consumer, upgrader: upgrader}
}

func (c *consumeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade websocket connection: " + err.Error() + "\n")
		return
	}
	defer connection.Close()

	for {
		if msg, err := c.consumer.consume(); err == nil {
			if err := connection.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("failed to write websocket message: " + err.Error() + "\n")
			}
		} else {
			log.Printf("failed to consume from topic: " + err.Error() + "\n")
		}
	}
}

type consumer struct {
	kc *kafka.Consumer
}

func newConsumer(bootstrapServers, username, password, topic string) (*consumer, error) {
	rand.Seed(time.Now().UnixNano())
	config := kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          "random-" + strconv.Itoa(int(rand.Uint32())),
	}
	if username != "" {
		config["security.protocol"] = "SASL_PLAINTEXT"
		config["sasl.mechanisms"] = "PLAIN"
		config["sasl.username"] = username
		config["sasl.password"] = password
	}

	kc, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Printf("failed to construct consumer: " + err.Error() + "\n")
		return nil, err
	}

	kc.SubscribeTopics([]string{topic}, nil)

	return &consumer{kc: kc}, nil
}

func (c *consumer) consume() ([]byte, error) {
	if message, err := c.kc.ReadMessage(-1); err != nil {
		log.Printf("failed to consume message from topic: " + err.Error() + "\n")
		return nil, err
	} else {
		return message.Value, nil
	}
}

func (c *consumer) close() {
	c.kc.Close()
}
