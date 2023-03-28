package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func index(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("got error while reading body: %s", err)
	}
	log.Info("received request", "from", req.RemoteAddr, "user agent", req.UserAgent(), "path", req.RequestURI, "method", req.Method, "body", "\n"+string(body))
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

func health(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

func liveness(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

func main() {
	port, port_set := os.LookupEnv("LISTEN_PORT")
	if !port_set {
		port = "8080"
	}

	address := os.Getenv("LISTEN_ADDRESS")
	listen_on := fmt.Sprintf("%s:%s", address, port)

	fmt.Println(`Welcome to this simple echo server 👋🏻
This is a simple HTTP server that will answer a 200 OK message to every request.
For every request that arrives the server will print the sender's address, the path used and the body's content.
You can change the listening address by setting the LISTEN_ADDRESS and LISTEN_PORT environment variables.

Available endpoints:
  - /metrics/  Prometheus metrics
  - /health/   returns always 200. Could change in the future
  - /liveness/ returns always 200. Could change in the future
  - /<path>/   returns 200. The request details get printed to stdout.`)
	fmt.Println("")

	http.HandleFunc("/", index)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", health)
	http.HandleFunc("/liveness", liveness)

	log.Info("starting server...", "address", address, "port", port)
	log.Fatal(http.ListenAndServe(listen_on, nil))
}
