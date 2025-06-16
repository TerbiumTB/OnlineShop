package mq

import (
	"context"
	"github.com/google/uuid"
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
	m kafka.Message
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
	return &Kafka{w, r, kafka.Message{}}
}

func (k *Kafka) Close() (err error) {
	err = k.w.Close()
	if err != nil {
		return
	}
	return k.r.Close()
}

func (k *Kafka) Send(e *model.Event) (err error) {
	err = k.w.WriteMessages(context.Background(), kafka.Message{
		Key:   e.ID[:],
		Value: e.Payload,
	})
	return
}

func (k *Kafka) Receive() (e *model.Event, err error) {
	k.m, err = k.r.FetchMessage(context.Background())
	//k.r.FetchMessage()
	if err != nil {
		return nil, err
	}
	id, err := uuid.FromBytes(k.m.Key)

	if err != nil {
		return nil, err
	}
	return model.NewEventWithID(id, k.m.Value), nil
}

func (k *Kafka) Register() error {
	return k.r.CommitMessages(context.Background(), k.m)
}
