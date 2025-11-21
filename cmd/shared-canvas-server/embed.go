package main

import (
	"embed"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

//go:embed web-dist/*
var embeddedWeb embed.FS

func createEmbedHandler() http.HandlerFunc {
	fs := http.FS(embeddedWeb)
	return func(w http.ResponseWriter, r *http.Request) {
		embedHandler(fs, w, r)
	}
}

func embedHandler(fs http.FileSystem, w http.ResponseWriter, r *http.Request) {
	// Clean the path and map to the embedded web-dist directory
	p := path.Clean(r.URL.Path)
	if p == "/" || p == "." {
		p = "/web-dist/index.html"
	} else {
		p = "/web-dist" + p
	}

	f, err := fs.Open(p)
	if err != nil {
		// Fallback to index.html for unknown routes (SPA support)
		if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/ws") {
			if idx, err2 := fs.Open("/web-dist/index.html"); err2 == nil {
				defer func(idx http.File) {
					err := idx.Close()
					if err != nil {
						log.Printf("Error closing file: %v", err)
					}
				}(idx)
				http.ServeContent(w, r, "index.html", time.Time{}, idx)
				return
			}
		}
		http.NotFound(w, r)
		return
	}
	defer func(f http.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}(f)
	http.ServeContent(w, r, path.Base(p), time.Time{}, f)
}
