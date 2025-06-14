package main

import (
	"context"
	"github.com/google/uuid"
	"orders/application"
	"orders/infrastructure/inoutbox"
	"orders/infrastructure/mq"
	"orders/infrastructure/storage"
	"orders/infrastructure/trx"
	"orders/pkg/postgres"
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

	db, err := postgres.Init()
	if err != nil {
		lg.Fatal(err)
	}
	s, err := storage.NewOrderDB(db)
	o, err := inoutbox.NewOutbox(db)
	m := trx.NewDBManager(db)

	if err != nil {
		lg.Fatal(err)
	}

	k := mq.NewKafka()
	defer k.Close()

	service := application.NewService(s, o, m)
	outbox := application.NewOutboxWorker(k, o, lg)

	outbox.Start(context.Background(), 3*time.Second)

	//order := models.Order{}
	for i := 0; i < 10; i++ {
		//order := models.NewOrder(uuid.New(), float64(i)*12+100, "Labubu")
		err = service.Add(uuid.New().String(), float64(i)*12+100, "Labubu")
		if err != nil {
			lg.Println(err)
		}
		time.Sleep(6 * time.Second)
	}
	//db, err := storage.NewOrderDB(nil)
	//application.Service{db}

	//lg.Println("Starting orders server")
	//
	//kafkaURL := os.Getenv("KAFKA_BROKER")
	//topicOrders := os.Getenv("KAFKA_ORDERS_TOPIC")
	//
	//p := kafka.Writer{
	//	Addr:     kafka.TCP(kafkaURL),
	//	Topic:    topicOrders,
	//	Balancer: &kafka.LeastBytes{},
	//}
	//defer p.Close()
	//
	//for i := 0; i < 100; i += 1 {
	//	msg := fmt.Sprintf("Order %d", i)
	//	err := p.WriteMessages(context.Background(), kafka.Message{
	//		Value: []byte(msg),
	//	})
	//	if err != nil {
	//		lg.Fatal(err)
	//	}
	//	time.Sleep(3 * time.Second)
	//}

}
