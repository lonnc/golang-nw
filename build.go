package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
)

func main() {
	w, err := os.Create("client.nw")
	if err != nil {
		panic(err)
	}
	defer w.Close()

	zw := zip.NewWriter(w)
	defer zw.Close()

	for _, f := range []string{"package.json", "index.html", "client.js", "script.js", "myapp.exe"} {
		fw, err := zw.Create(f)
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}
		_, err = fw.Write(b)
		if err != nil {
			panic(err)
		}
	}
}
