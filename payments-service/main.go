package main

import (
	"context"
	"log"
	"os"
	"payments/application"
	"payments/infrastructure/inoutbox"
	"payments/infrastructure/mq"
	"payments/pkg/postgres"
	"time"
)

func main() {
	lg := log.New(os.Stdout, "Payments ", log.LstdFlags)

	lg.Println("Starting payments server")

	db, err := postgres.Init()
	if err != nil {
		lg.Fatal(err)
	}
	//s, err := storage.NewOrderDB(db)
	i, err := inoutbox.NewInbox(db)

	if err != nil {
		lg.Fatal(err)
	}

	k := mq.NewKafka()
	defer k.Close()

	inbox := application.NewInboxWorker(k, i, lg)

	inbox.Start(context.Background(), 2*time.Second)

	for {
	}
}
