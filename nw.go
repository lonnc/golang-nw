package nw

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

const EnvVar = "GOLANG-NW"

var ErrMissingEnvVariable = errors.New("missing environment variable '" + EnvVar + "'")

type NodeWebkit struct {
	Url string // URL to issue callback command to
}

func New() (NodeWebkit, error) {
	url := os.Getenv(EnvVar)
	if url == "" {
		return NodeWebkit{}, ErrMissingEnvVariable
	}
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	return NodeWebkit{url}, nil
}

// ListenAndServe listens to a random port on localhost
// and the issues a Redirect to node-webkit.
// If handler is nil then the http.DefaultHandler will be used
func (n NodeWebkit) ListenAndServe(handler http.Handler) error {
	httpAddr := "127.0.0.1:0"
	listener, err := net.Listen("tcp", httpAddr)
	// If we were not able to establish a socket send error back to node-webkit
	if err != nil {
		// Send error message to node-webkit
		return n.Error(err.Error())
	}
	defer listener.Close()

	errs := make(chan error)
	defer close(errs)

	// Issue redirect asynchronously so we can get on with the business of serving http requests
	go func() {
		if err := n.Redirect("http://" + listener.Addr().String() + "/"); err != nil {
			errs <- err
		}
	}()

	// Now start the normal http listener
	go func() {
		if err := http.Serve(listener, handler); err != nil {
			errs <- err
		}
	}()

	// Wait for an error from either serving or the redirect
	if err := <-errs; err != nil {
		// and forward it over to node-webkit
		return n.Error(err.Error())
	}
	return nil
}

// Redirect sends a redirect message to node-webkit
func (n NodeWebkit) Redirect(url string) error {
	return n.send("redirect", url)
}

// Redirect sends an error message to node-webkit
func (n NodeWebkit) Error(msg string) error {
	return n.send("error", msg)
}

func (n NodeWebkit) send(key string, value string) error {
	r, err := http.Post(n.Url+key, "text/plain", strings.NewReader(value))
	if err != nil {
		return err
	}
	if r.StatusCode < 200 || r.StatusCode > 299 {
		return fmt.Errorf("Unexpected status code: %d", r.StatusCode)
	}
	return nil
}
