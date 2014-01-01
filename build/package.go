package build

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/lonnc/golang-nw"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

type Package struct {
	Name   string `json:"name"`
	Main   string `json:"main"`
	Window Window `json:"window"`
	Bin    string `json:"-"`
	EnvVar string `json:"-"`
}

type Window struct {
	Title    string `json:"title,omitempty"`
	Toolbar  bool   `json:"toolbar"`
	Show     bool   `json:"show,omitempty"`
	Position string `json:"position,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type Templates struct {
	IndexHtml string
	ClientJs  string
	ScriptJs  string
}

var DefaultTemplates = Templates{IndexHtml: index, ClientJs: client, ScriptJs: script}

// CreateNW creates a node-webkit .nw file
func (p Package) CreateNW(zw *zip.Writer, templates Templates, myapp io.Reader, includes string) error {
	// Add in a couple of package defaults
	p.Main = "index.html"
	p.EnvVar = nw.EnvVar

	if w, err := zw.Create("package.json"); err != nil {
		return err
	} else {
		if _, err := p.writeJsonTo(w); err != nil {
			return err
		}
	}

	filenameTemplates := map[string]string{
		"index.html": templates.IndexHtml,
		"client.js":  templates.ClientJs,
		"script.js":  templates.ScriptJs}
	for filename, str := range filenameTemplates {
		if w, err := zw.Create(filename); err != nil {
			return err
		} else {
			if t, err := template.New(filename).Parse(str); err != nil {
				return err
			} else {
				if err := t.Execute(w, p); err != nil {
					return err
				}
			}
		}
	}

	if includes != "" {
		if err := copyIncludes(zw, includes); err != nil {
			return err
		}
	}
    
	binHeader := zip.FileHeader{Name: p.Bin}
	binHeader.SetMode(0755) // Make it executable
	if w, err := zw.CreateHeader(&binHeader); err != nil {
		return err
	} else {
		if _, err := io.Copy(w, myapp); err != nil {
			return err
		}
	}

	return nil
}

func (p Package) writeJsonTo(w io.Writer) (int64, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

// Copy any files from the includes directory
func copyIncludes(zw *zip.Writer, includes string) (includeErr error) {
	includes = path.Clean(includes)
	if !strings.HasSuffix(includes, "/") {
		includes += "/"
	}
	filepath.Walk(includes, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			includeErr = err
			return err
		}
		if info.IsDir() {
			return nil
		}
		p := strings.TrimPrefix(path, includes)
		fmt.Printf("Path: %s\nPrefix: %s\n", path, includes)
		fmt.Printf("Adding %s to the zip file\n", p)
		if w, err := zw.Create(p); err != nil {
			includeErr = err
			return err
		} else {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				includeErr = err
				return err
			} else {
				_, err := w.Write(b)
				if err != nil {
					includeErr = err
					return err
				}
			}
		}
		return nil
	})
	return
}

