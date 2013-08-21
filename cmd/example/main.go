package main

import (
	"fmt"
	"github.com/lonnc/golang-nw"
	"net/http"
)

func main() {
	// Setup our handler
	http.HandleFunc("/", hello)

	// Create a link back to node-webkit using the environment variable
	// populated by golang-nw's node-webkit code
	nodeWebkit, err := nw.New()
	if err != nil {
		panic(err)
	}

	// Pick a random localhost port, start listening for http requests
	// and send a message address back to node-webkit to redirect
	if err := nodeWebkit.ListenAndServe(nil); err != nil {
		panic(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from golang.")
}
