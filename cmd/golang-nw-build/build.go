package main

import (
	"archive/zip"
	"flag"
	"github.com/lonnc/golang-nw/build"
	"os"
	"path/filepath"
)

var (
	app  = "myapp.exe"
	out  = "myapp.nw"
	name = "My Application"
)

func main() {
	flag.StringVar(&app, "app", app, "Application to be wrapped by node-webkit.")
	flag.StringVar(&name, "name", name, "Application name.")
	flag.StringVar(&out, "out", out, "Destination file for generated node-webkit .nw file.")
	flag.Parse()

	w, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	zw := zip.NewWriter(w)
	defer zw.Close()

	r, err := os.Open(app)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	bin := filepath.Base(app)
	p := build.Package{Name: name, Bin: bin, Window: build.Window{Title: name}}

	if err := p.CreateNW(zw, build.DefaultTemplates, r); err != nil {
		panic(err)
	}
}
