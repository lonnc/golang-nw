/*
Call a golang web application from node-webkit to get a native looking application.


Instructions


Go get golang-nw:

    go get github.com/lonnc/golang-nw/cmd/golang-nw-pkg

Create an app:

See https://github.com/lonnc/golang-nw/blob/master/cmd/example/main.go
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
		// and send a message back to node-webkit to redirect
		if err := nodeWebkit.ListenAndServe(nil); err != nil {
			panic(err)
		}
	}

	func hello(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from golang.")
	}


Build your app:

    go install .\src\github.com\lonnc\golang-nw\cmd\example


Wrap it in node-webkit:

    .\bin\golang-nw-pkg.exe -app=.\bin\example.exe -name="My Application" -out="myapp.exe"


Finally execute node-webkit with the myapp.nw generated above as a parameter:

    nw.exe myapp.nw


You will probably want to create your own build script so you can control toolbar visibility,
window dimensions etc.
See http://godoc.org/github.com/lonnc/golang-nw/build and
https://github.com/lonnc/golang-nw/blob/master/cmd/golang-nw-pkg/pkg.go
*/
package nw
