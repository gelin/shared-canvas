package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

//go:embed web/*
var embeddedWeb embed.FS

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.IntVar(&port, "p", 8080, "Port to listen on (shorthand)")
	flag.Parse()

	addr := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":  "ok",
			"service": "shared-canvas",
			"time":    time.Now().Format(time.RFC3339Nano),
		})
	})

	mux.HandleFunc("/api/time", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"now":   time.Now().Format(time.RFC3339Nano),
			"epoch": time.Now().Unix(),
		})
	})

	// Static file server from embedded assets
	fs := http.FS(embeddedWeb)
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean path and map to embedded web directory
		p := path.Clean(r.URL.Path)
		if p == "/" || p == "." {
			p = "/web/index.html"
		} else {
			p = "/web" + p
		}

		f, err := fs.Open(p)
		if err != nil {
			// Fallback to index.html for unknown routes (SPA support)
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				if idx, err2 := fs.Open("/web/index.html"); err2 == nil {
					defer idx.Close()
					http.ServeContent(w, r, "index.html", time.Time{}, idx)
					return
				}
			}
			http.NotFound(w, r)
			return
		}
		defer f.Close()
		http.ServeContent(w, r, path.Base(p), time.Time{}, f)
	}))

	srv := &http.Server{Addr: addr, Handler: loggingMiddleware(mux)}

	// Graceful shutdown
	go func() {
		log.Printf("Server listening on http://localhost%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)
		dur := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.status, dur)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
