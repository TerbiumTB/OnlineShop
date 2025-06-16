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
	sm := mux.NewRouter()
	sm.PathPrefix("/order/").Handler(handlers.NewOrdersHandler())

	server.StartServer(":8000", sm, lg)
	//uploadRouter := sm.Methods(http.MethodPost).Subrouter()
	//uploadRouter.HandleFunc("/upload/{filename}", handlers.StorageHandler())
}
