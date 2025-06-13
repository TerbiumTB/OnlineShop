package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	//"github.com/gorilla/mux"
	"log"
	"os"
	//	"os/signal"
	//	"time"
)

func main() {
	lg := log.New(os.Stdout, "Payments ", log.LstdFlags)

	lg.Println("Starting payments server")

	kafkaURL := os.Getenv("KAFKA_BROKER")
	topicOrders := os.Getenv("KAFKA_ORDERS_TOPIC")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topicOrders,
		GroupID:  "payments",
		MinBytes: 10e2,
		MaxBytes: 10e6,
	})
	defer r.Close()

	lg.Printf("Reading from kafka")

	for {
		msg, err := r.ReadMessage(context.Background())

		if err != nil {
			log.Printf("Error reading message: %v", err)
		}

		log.Printf("Message received: %s", msg.Value)

	}
	//server.StartServer(":8080", http.DefaultServeMux, lg)
}
