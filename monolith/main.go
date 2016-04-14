package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/braintree/manners"
	"github.com/kelseyhightower/app/handlers"
	"github.com/kelseyhightower/app/health"
	"github.com/kelseyhightower/app/user"
)

func main() {
	log.Println("Starting server on 0.0.0.0:8080...")
	log.Println("Listening on 0.0.0.0:8000 for health checks")

	flag.Parse()
	errChan := make(chan error, 10)

	hmux := http.NewServeMux()
	hmux.HandleFunc("/healthz", health.HealthzHandler)
	hmux.HandleFunc("/readiness", health.ReadinessHandler)
	hmux.HandleFunc("/healthz/status", health.HealthzStatusHandler)
	hmux.HandleFunc("/readiness/status", health.ReadinessStatusHandler)
	hserver := manners.NewServer()
	hserver.Addr = ":8000"
	hserver.Handler = handlers.LoggingHandler(hmux)

	go func() {
		errChan <- hserver.ListenAndServe()
	}()

	mux := http.NewServeMux()
	mux.Handle("/login", handlers.LoginHandler("123456789", user.DB))
	mux.Handle("/", handlers.JWTAuthHandler(handlers.HelloHandler))
	server := manners.NewServer()
	server.Addr = ":8080"
	server.Handler = handlers.LoggingHandler(mux)

	go func() {
		errChan <- server.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Started successfully.")
	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			health.SetReadinessStatus(http.StatusServiceUnavailable)
			server.BlockingClose()
			os.Exit(0)
		}
	}
}
