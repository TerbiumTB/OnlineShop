package main

import (
	"context"
	"github.com/gorilla/mux"
	swag "github.com/swaggo/http-swagger"
	"net/http"
	_ "orders/docs"
	"orders/handlers"
	"orders/infrastructure/inoutbox"
	"orders/infrastructure/mq"
	"orders/infrastructure/storage"
	"orders/infrastructure/trx"
	"orders/pkg/postgres"
	"orders/pkg/server"
	"orders/services"
	"time"
	//"github.com/gorilla/mux"
	"log"
	"os"
)

// @title Orders
// @version 1.0

// @host localhost:8080
// @BasePath /
func main() {
	lg := log.New(os.Stdout, "Orders ", log.LstdFlags)

	db, err := postgres.Init()
	if err != nil {
		lg.Fatal(err)
	}

	orders, err := storage.NewOrderDB(db)
	outbox, err := inoutbox.NewOutbox(db)
	trxManager := trx.NewDBManager(db)

	if err != nil {
		lg.Fatal(err)
	}

	broker := mq.NewKafka()
	defer broker.Close()

	service.NewOutboxWorker(broker, outbox, lg).Start(context.Background(), 2*time.Second)
	service.NewStatusWorker(broker, orders, lg).Start(context.Background(), 2*time.Second)

	s := service.NewOrderService(orders, outbox, trxManager)

	h := handlers.NewHandler(s, lg)
	sm := mux.NewRouter()

	createRouter := sm.Methods(http.MethodPost).Subrouter()
	createRouter.HandleFunc("/order/create/{user_id}", h.CreateOrder)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/order/{id}", h.GetOrder)

	allRouter := sm.Methods(http.MethodGet).Subrouter()
	allRouter.HandleFunc("/order", h.AllOrders)

	sm.PathPrefix("/docs/").Handler(swag.WrapHandler)

	server.StartServer(":8080", sm, lg)
}
