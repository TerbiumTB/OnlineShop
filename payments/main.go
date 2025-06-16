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
// @BasePath /
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
	service.NewPaymentWorker(accounts, inbox, outbox, manager, lg).StartPaying(context.Background(), 2*time.Second)

	as := service.NewAccountService(accounts, lg)

	lg.Println("tuning server...")

	h := handler.NewHandler(as, lg)
	sm := mux.NewRouter()

	createRouter := sm.Methods(http.MethodPost).Subrouter()
	createRouter.HandleFunc("/account/create/{user_id}", h.CreateAccount)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/account/{user_id}", h.GetAccount)

	allRouter := sm.Methods(http.MethodGet).Subrouter()
	allRouter.HandleFunc("/account", h.AllAccounts)

	updateRouter := sm.Methods(http.MethodPatch).Subrouter()
	updateRouter.HandleFunc("/account/update/{user_id}", h.UpdateBalance)

	sm.PathPrefix("/docs/").Handler(swag.WrapHandler)

	server.StartServer(":8080", sm, lg)

}
