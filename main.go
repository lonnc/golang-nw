package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	httpAddr = "localhost:8080"
)

func main() {
	flag.StringVar(&httpAddr, "http", httpAddr, "HTTP service address (e.g., 'localhost:8080')")
	flag.Parse()

	l := log.New(os.Stdout, "", log.Flags())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from golang.")
	})

	if ln, err := net.Listen("tcp", httpAddr); err != nil {
		l.Fatalf("Failed to start listener on %s: %s", httpAddr, err)
	} else {
		l.Printf("HTTP listening on %s", ln.Addr().String())
		if err := http.Serve(ln, nil); err != nil {
			l.Fatalf("Failed to server http: %s", err)
		}
	}
}
