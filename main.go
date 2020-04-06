package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// TODO logging
// TODO better error handling

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		log.Fatalf("HOST environment variable must be set.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT environment variable must be set.")
	}

	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrapServers == "" {
		log.Fatalf("KAFKA_BOOTSTRAP_SERVERS environment variable must be set.")
	}

	kafkaUsername := os.Getenv("KAFKA_USERNAME")
	kafkaPassword := os.Getenv("KAFKA_PASSWORD")
	if (kafkaUsername == "") != (kafkaPassword == "") {
		log.Fatalf("either provide both KAFKA_USERNAME and KAFKA_PASSWORD for Kafka with SASL PLAIN authentication\n" +
			"  or provide neither for Kafka with no authentication")
	}

	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	if kafkaTopic == "" {
		log.Fatalf("KAFKA_TOPIC environment variable must be set.")
	}

	consumer, err := newConsumer(kafkaBootstrapServers, kafkaUsername, kafkaPassword, kafkaTopic)
	if err != nil {
		log.Fatalf("failed constructing Kafka consumer: " + err.Error())
	}
	defer consumer.close()

	producer, err := newProducer(kafkaBootstrapServers, kafkaUsername, kafkaPassword, kafkaTopic)
	if err != nil {
		log.Fatalf("failed constructing Kafka consumer: " + err.Error())
	}
	defer producer.close()

	upgrader := new(websocket.Upgrader)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "produce.html") })
	http.Handle("/produce.ws", newProduceHandler(producer, upgrader))
	http.HandleFunc("/consume", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "consume.html") })
	http.Handle("/consume.ws", newConsumeHandler(consumer, upgrader))

	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("failed to listen-and-serve: " + err.Error())
	}
}
