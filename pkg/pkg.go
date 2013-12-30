package pkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Pkg struct {
	Url          string
	Bin          string
	Dependencies []string
}

func New(version string, nwOs string, nwArch string) Pkg {
	pkgOs, ok := pkgOss[nwOs]
	if !ok {
		panic(fmt.Errorf("Unsupported os %q", nwOs))
	}

	var arch string
	switch nwArch {
	case "386":
		arch = "ia32"
	case "amd64":
		arch = "x64"
	default:
		panic(fmt.Errorf("Unsupported arch %q", nwArch))
	}

	url := fmt.Sprintf("https://s3.amazonaws.com/node-webkit/%s/node-webkit-%s-%s-%s%s", version, version, pkgOs.os, arch, pkgOs.ext)

	pkg := Pkg{
		Url:          url,
		Bin:          pkgOs.bin,
		Dependencies: pkgOs.deps,
	}

	return pkg
}

// Download retrieves the Url into the passed directory.
// If destDir=="" then TempDir is used.
// The file path is returned upon completion
func (p Pkg) Download(destDir string) (string, error) {
	if destDir == "" {
		destDir = os.TempDir()
	}

	// Where do we want to download to
	out := filepath.Join(destDir, path.Base(p.Url))

	// See if we already have it
	if exists, err := isExists(out); err != nil {
		return out, err
	} else if exists {
		return out, nil
	}

	// Download into memory then write to disk after
	client := http.DefaultClient
	r, err := client.Get(p.Url)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return "", fmt.Errorf("Failed to download %q got: %s", p.Url, r.Status)
	}

	if content, err := ioutil.ReadAll(r.Body); err != nil {
		return "", err
	} else {
		if err := ioutil.WriteFile(out, content, 0666); err != nil {
			return "", err
		}
	}

	return out, nil
}

// Package wraps populates destDir with the node-webkit depedencies and cat nw.exe [nw content] > binName
func (p Pkg) Package(nodeWebkitPath string, nw io.Reader, binName string, destDir string) error {
	// Ensure destDir exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	// Check we have a zip
	if nodeWebkitPathZip, err := ensureZip(nodeWebkitPath); err != nil {
		return err
	} else {
		nodeWebkitPath = nodeWebkitPathZip
	}

	// Extract dependencies from zip file
	zr, err := zip.OpenReader(nodeWebkitPath)
	if err != nil {
		return err
	}
	defer zr.Close()

	// Get list of files in the zip archive, excluding the preceding directory
	zipFiles := map[string]*zip.File{}
	for _, f := range zr.File {
		zipFiles[path.Base(f.Name)] = f
	}

	if p.Bin != "" {
		if bin, ok := zipFiles[p.Bin]; !ok {
			return fmt.Errorf("Failed to find %s in %s", p.Bin, nodeWebkitPath)
		} else {
			if err := p.copyBin(bin, nw, binName, destDir); err != nil {
				return err
			}
		}
	}

	if err := p.copyDependencies(zipFiles, destDir); err != nil {
		return err
	}

	return nil
}

func (p Pkg) copyBin(bin *zip.File, nw io.Reader, binName string, destDir string) error {
	r, err := bin.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	filename := filepath.Join(destDir, binName)
	w, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer w.Close()

	// Copy nw binary
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	// Copy nw
	_, err = io.Copy(w, nw)
	if err != nil {
		return err
	}

	return nil
}

func (p Pkg) copyDependencies(zipFiles map[string]*zip.File, destDir string) error {
	// And extract the dependencies
	for _, dep := range p.Dependencies {
		filename := filepath.Join(destDir, dep)

		// Only copy over if it doesn't already exist
		if exists, err := isExists(filename); err != nil {
			return err
		} else if exists {
			continue
		}

		// And copy it over
		var r io.ReadCloser = nil
		if zipFile, ok := zipFiles[dep]; !ok {
			return fmt.Errorf("Failed to find %s", dep)
		} else {
			if f, err := zipFile.Open(); err != nil {
				return err
			} else {
				r = f
			}
		}
		defer r.Close()

		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()

		if _, err = io.Copy(w, r); err != nil {
			return err
		}

	}

	return nil
}

type pkgOs struct {
	os   string
	bin  string
	deps []string
	ext  string
}

var windows = pkgOs{
	os:   "win",
	bin:  "nw.exe",
	deps: []string{"ffmpegsumo.dll", "icudt.dll", "libEGL.dll", "libGLESv2.dll", "nw.pak"},
	ext:  ".zip",
}

var linux = pkgOs{
	os:   "linux",
	bin:  "nw",
	deps: []string{"libffmpegsumo.so", "nw.pak"},
	ext:  ".tar.gz",
}

var darwin = pkgOs{
	os:   "osx",
	bin:  "",
	deps: []string{},
	ext:  ".zip",
}

var pkgOss = map[string]pkgOs{
	"windows": windows,
	"linux":   linux,
	"darwin":  darwin,
}

func isExists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Check whether this is a .zip and return if not create a zip file and return that
func ensureZip(filename string) (string, error) {
	if strings.HasSuffix(filename, ".zip") {
		return filename, nil
	}

	if strings.HasSuffix(filename, ".tar.gz") {
		filenameZip := filename[:len(filename)-len(".tar.gz")] + ".zip"
		if exists, err := isExists(filenameZip); err != nil {
			return "", err
		} else {
			if !exists {
				if err := toZip(filename, filenameZip); err != nil {
					return "", err
				}
			}
			return filenameZip, err
		}
	}

	return "", fmt.Errorf("Do not know how to create a zip archive from %s", filename)
}

// convert a .tar.gz into a .zip
func toZip(filenameTarGz string, filenameZip string) error {
	r, err := os.Open(filenameTarGz)
	if err != nil {
		return err
	}
	defer r.Close()

	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}

	tgz := tar.NewReader(gz)

	filenameZipTmp := filenameZip + ".tmp"
	w, err := os.Create(filenameZipTmp)
	if err != nil {
		return err
	}
	defer func() {
		if w != nil {
			w.Close()
		}
	}()

	z := zip.NewWriter(w)
	defer func() {
		if z != nil {
			z.Close()
		}
	}()

	for {
		hdr, err := tgz.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		f, err := z.Create(hdr.Name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, tgz); err != nil {
			return err
		}
	}

	// Ok all done, lets close
	if err := z.Close(); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	// And rename
	if err := os.Rename(filenameZipTmp, filenameZip); err != nil {
		return err
	}
	return nil
}
