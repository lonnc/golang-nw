package main

import (
	"archive/zip"
	"flag"
	"github.com/lonnc/golang-nw/build"
	"os"
	"path/filepath"
)

var (
	app         = "myapp.exe"
	out         = "myapp.nw"
	name        = "My Application"
	includesDir = ""
)

func main() {
	flag.StringVar(&app, "app", app, "Application to be wrapped by node-webkit.")
	flag.StringVar(&name, "name", name, "Application name.")
	flag.StringVar(&out, "out", out, "Destination file for generated node-webkit .nw file.")
	flag.StringVar(&includesDir, "includesDir", includesDir, "Directory containing additional files to bundle with the .nw file")
	flag.Parse()

	if err := nwBuild(); err != nil {
		panic(err)
	}
}

func nwBuild() error {
	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()

	zw := zip.NewWriter(w)
	defer zw.Close()

	r, err := os.Open(app)
	if err != nil {
		return err
	}
	defer r.Close()

	bin := filepath.Base(app)
	p := build.Package{Name: name, Bin: bin, Window: build.Window{Title: name}}

	if err := p.CreateNW(zw, build.DefaultTemplates, r, includesDir); err != nil {
		return err
	}

	return nil
}
