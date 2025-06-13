package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"time"

	//"github.com/gorilla/mux"
	"log"
	"os"
)

type Event struct {
	Id   uuid.UUID
	Type string
	Msg  string
}

func main() {
	lg := log.New(os.Stdout, "Orders ", log.LstdFlags)

	lg.Println("Starting orders server")

	kafkaURL := os.Getenv("KAFKA_BROKER")
	topicOrders := os.Getenv("KAFKA_ORDERS_TOPIC")

	lg.Println(os.Environ())
	p := kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topicOrders,
		Balancer: &kafka.LeastBytes{},
	}
	defer p.Close()

	for i := 0; i < 100; i += 1 {
		msg := fmt.Sprintf("Order %d", i)
		err := p.WriteMessages(context.Background(), kafka.Message{
			Value: []byte(msg),
		})
		if err != nil {
			lg.Fatal(err)
		}
		time.Sleep(3 * time.Second)
	}

}
