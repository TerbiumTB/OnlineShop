package main

import (
	"apigateway/handlers"
	"github.com/gorilla/mux"
	"log"
	//"net/http"
	"apigateway/pkg/server"
	"os"
)

func main() {
	lg := log.New(os.Stdout, "API gateway ", log.LstdFlags)
	r := mux.NewRouter()
	r.PathPrefix("/order/").Handler(handlers.NewOrdersHandler())
	r.PathPrefix("/payment/").Handler(handlers.NewPaymentsHandler())

	server.StartServer(":8000", r, lg)
}
