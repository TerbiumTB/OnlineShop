package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartServer(address string, handler http.Handler, lg *log.Logger) {
	server := http.Server{
		Addr:         address,
		Handler:      handler,
		ErrorLog:     lg,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		lg.Printf("Starting server on port %s...", address)

		err := server.ListenAndServe()
		if err != nil {
			lg.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, os.Kill)

	sig := <-ch
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
