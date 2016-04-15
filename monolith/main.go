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
		httpAddr   = flag.String("http", "0.0.0.0:5000", "HTTP service address.")
		healthAddr = flag.String("health", "0.0.0.0:5001", "Health service address.")
		secret     = flag.String("secret", "secret", "JWT signing secret.")
		certFile   = flag.String("cert", "server.pem", "TLS certificate")
		keyFile    = flag.String("key", "server-key.pem", "TLS private key")
	)

	log.Println("Starting server...")
	log.Printf("HTTP service listening on %s", *httpAddr)
	log.Printf("Health service listening on %s", *healthAddr)

	flag.Parse()
	errChan := make(chan error, 10)

	hmux := http.NewServeMux()
	hmux.HandleFunc("/healthz", health.HealthzHandler)
	hmux.HandleFunc("/readiness", health.ReadinessHandler)
	hmux.HandleFunc("/healthz/status", health.HealthzStatusHandler)
	hmux.HandleFunc("/readiness/status", health.ReadinessStatusHandler)
	hserver := manners.NewServer()
	hserver.Addr = *healthAddr
	hserver.Handler = handlers.LoggingHandler(hmux)

	go func() {
		errChan <- hserver.ListenAndServe()
	}()

	mux := http.NewServeMux()
	mux.Handle("/login", handlers.LoginHandler(*secret, user.DB))
	mux.Handle("/", handlers.JWTAuthHandler(handlers.HelloHandler))
	server := manners.NewServer()
	server.Addr = *httpAddr
	server.Handler = handlers.LoggingHandler(mux)

	go func() {
		errChan <- server.ListenAndServeTLS(*certFile, *keyFile)
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
