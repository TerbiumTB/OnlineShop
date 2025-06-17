package main

import (
	"context"
	"github.com/gorilla/mux"
	swag "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	_ "payments/docs"
	"payments/handlers"
	"payments/infrastructure/inoutbox"
	"payments/infrastructure/mq"
	"payments/infrastructure/storage"
	"payments/infrastructure/trx"
	"payments/pkg/postgres"
	"payments/pkg/server"
	"payments/services"
	"time"
)

// @title Payments
// @version 1.0

// @host localhost:8081
// @BasePath /payment/
func main() {
	lg := log.New(os.Stdout, "Payments ", log.LstdFlags)

	lg.Println("starting db...")
	db, err := postgres.Init()
	if err != nil {
		lg.Fatal(err)
	}
	accounts, err := storage.NewAccountDB(db)
	if err != nil {
		lg.Fatal(err)
	}

	lg.Println("starting inbox...")
	inbox, err := inoutbox.NewInbox(db)
	if err != nil {
		lg.Fatal(err)
	}

	lg.Println("starting outbox...")
	outbox, err := inoutbox.NewOutbox(db)
	if err != nil {
		lg.Fatal(err)
	}
	lg.Println("starting manager...")
	manager := trx.NewDBManager(db)

	lg.Println("starting broker...")
	broker := mq.NewKafka()
	defer broker.Close()

	lg.Println("starting workers...")
	service.NewInboxWorker(broker, inbox, lg).Start(context.Background(), 2*time.Second)
	service.NewOutboxWorker(broker, outbox, lg).Start(context.Background(), 2*time.Second)
	service.NewPaymentWorker(accounts, inbox, outbox, manager, lg).StartPaying(context.Background(), 2*time.Second)

	s := service.NewAccountService(accounts, lg)

	lg.Println("tuning server...")

	h := handler.NewHandler(s, lg)
	r := mux.NewRouter()

	sr := r.PathPrefix("/payment/account/").Subrouter()
	sr.Methods(http.MethodPost).Path("/create/{user_id}").HandlerFunc(h.CreateAccount)
	sr.Methods(http.MethodGet).Path("/get/{user_id}").HandlerFunc(h.GetAccount)
	sr.Methods(http.MethodGet).Path("/get").HandlerFunc(h.AllAccounts)
	sr.Methods(http.MethodPatch).Path("/update/{user_id}").HandlerFunc(h.UpdateBalance)

	r.PathPrefix("/docs/").Handler(swag.WrapHandler)

	server.StartServer(":8080", r, lg)

}
