package mq

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"orders/models"
	"os"
)

var (
	kafkaURL      = os.Getenv("KAFKA_BROKER")
	topicOrders   = os.Getenv("KAFKA_ORDERS_TOPIC")
	topicPayments = os.Getenv("KAFKA_PAYMENTS_TOPIC")
	groupID       = os.Getenv("KAFKA_GROUP_ID")
)

type Kafka struct {
	w *kafka.Writer
	r *kafka.Reader
}

func NewKafka() *Kafka {
	w := &kafka.Writer{
		Addr:  kafka.TCP(kafkaURL),
		Topic: topicOrders,

		Balancer: &kafka.LeastBytes{},
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topicPayments,
		GroupID:  groupID,
		MinBytes: 10e2,
		MaxBytes: 10e6,
	})

	//r := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{kafkaURL}})
	return &Kafka{w, r}
}

func (k *Kafka) Close() (err error) {
	err = k.w.Close()
	if err != nil {
		return
	}
	return k.r.Close()
}

func (k *Kafka) Send(e *models.Event) (err error) {
	err = k.w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(e.Type),
		Value: e.Payload,
	})
	return
}

func (k *Kafka) Receive(e *models.Event) (err error) {
	return fmt.Errorf("not implemented")
}
