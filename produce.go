package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type produceHandler struct {
	producer *producer
	upgrader *websocket.Upgrader
}

func newProduceHandler(producer *producer, upgrader *websocket.Upgrader) *produceHandler {
	return &produceHandler{producer: producer, upgrader: upgrader}
}

func (p *produceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	connection, err := p.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade websocket connection: " + err.Error() + "\n")
		return
	}
	defer connection.Close()

	for {
		if _, message, err := connection.ReadMessage(); err != nil {
			log.Printf("failed to read message from websocket connection: " + err.Error() + "\n")
			break
		} else {
			p.producer.produce(message)
		}
	}
}

type producer struct {
	kp    *kafka.Producer
	topic string
}

func newProducer(bootstrapServers, username, password, topic string) (*producer, error) {
	config := kafka.ConfigMap{"bootstrap.servers": bootstrapServers}
	if username != "" {
		config["security.protocol"] = "SASL_PLAINTEXT"
		config["sasl.mechanisms"] = "PLAIN"
		config["sasl.username"] = username
		config["sasl.password"] = password
	}

	kp, err := kafka.NewProducer(&config)
	if err != nil {
		log.Printf("failed to construct producer: " + err.Error() + "\n")
		return nil, err
	}

	return &producer{
		kp:    kp,
		topic: topic,
	}, nil
}

func (p *producer) produce(message []byte) {
	p.kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, nil)
	p.kp.Flush(0)
}

func (p *producer) close() {
	p.kp.Close()
}
