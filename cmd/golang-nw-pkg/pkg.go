package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"github.com/lonnc/golang-nw/build"
	"github.com/lonnc/golang-nw/pkg"
	"os"
	"path/filepath"
)

var (
	app      = "myapp.exe"
	name     = "My Application"
	bin      = "myapp.exe"
	binDir   = "."
	cacheDir = "."
)

func main() {
	flag.StringVar(&app, "app", app, "Application to be wrapped by node-webkit.")
	flag.StringVar(&name, "name", name, "Application name.")
	flag.StringVar(&bin, "bin", bin, "Destination file for combined application and node-webkit .nw file (will be placed in binDir directory).")
	flag.StringVar(&binDir, "binDir", binDir, "Destination directory for bin and dependencies.")
	flag.StringVar(&cacheDir, "cacheDir", cacheDir, "Directory to cache node-webkit download.")
	flag.Parse()

	p := pkg.Win32

	nw := filepath.Join(cacheDir, bin+".nw")
	fmt.Printf("Building:\t %s\n", nw)
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		panic(err)
	}
	if err := nwBuild(nw); err != nil {
		panic(err)
	}

	fmt.Printf("Downloading:\t %s\n", p.Url)
	nodeWebkitPath, err := nwDownload(p)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Packaging:\t %s\n", filepath.Join(binDir, bin))
	if err := nwPkg(p, nodeWebkitPath, nw); err != nil {
		panic(err)
	}
}

func nwBuild(nw string) error {
	w, err := os.Create(nw)
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

	if err := p.CreateNW(zw, build.DefaultTemplates, r); err != nil {
		return err
	}

	return nil
}

func nwDownload(p pkg.Pkg) (string, error) {
	return p.Download(cacheDir)
}

func nwPkg(p pkg.Pkg, nodeWebkitPath string, nw string) error {
	r, err := os.Open(nw)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := p.Package(nodeWebkitPath, r, bin, binDir); err != nil {
		return err
	}

	return nil
}
