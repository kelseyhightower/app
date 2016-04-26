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
	var (
		healthAddr = flag.String("health", "127.0.0.1:10081", "Health service address.")
		httpAddr   = flag.String("http", "0.0.0.0:10080", "HTTP service address.")
		secret     = flag.String("secret", "secret", "JWT signing secret.")
		certFile   = flag.String("cert", "server.pem", "TLS certificate.")
		keyFile    = flag.String("key", "server-key.pem", "TLS private key.")
	)

	log.Println("Starting server...")
	log.Printf("Health service listening on %s", *healthAddr)
	log.Printf("HTTP service listening on %s", *httpAddr)

	flag.Parse()
	errChan := make(chan error, 10)

	hmux := http.NewServeMux()
	hmux.HandleFunc("/healthz", health.HealthzHandler)
	hmux.HandleFunc("/readiness", health.ReadinessHandler)
	hmux.HandleFunc("/healthz/status", health.HealthzStatusHandler)
	hmux.HandleFunc("/readiness/status", health.ReadinessStatusHandler)
	healthServer := manners.NewServer()
	healthServer.Addr = *healthAddr
	healthServer.Handler = handlers.LoggingHandler(hmux)

	go func() {
		errChan <- healthServer.ListenAndServe()
	}()

	mux := http.NewServeMux()
	mux.Handle("/login", handlers.LoginHandler(*secret, user.DB))
	mux.Handle("/", handlers.JWTAuthHandler(handlers.HelloHandler))

	httpServer := manners.NewServer()
	httpServer.Addr = *httpAddr
	httpServer.Handler = handlers.LoggingHandler(mux)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			health.SetReadinessStatus(http.StatusServiceUnavailable)
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}
}
