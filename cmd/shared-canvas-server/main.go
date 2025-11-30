package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var imgHolder *ImageHolder

func main() {
	var port int
	var width int
	var height int
	flag.IntVar(&port, "port", 8080, "Port to listen on")
	flag.IntVar(&width, "width", 576, "Canvas width in pixels (576 by default)")
	flag.IntVar(&height, "height", 576, "Canvas height in pixels (576 by default)")
	flag.Parse()

	imgHolder = NewImageHolder(width, height)
	log.Printf("Canvas size set to %dx%d pixels\n", width, height)

	addr := fmt.Sprintf(":%d", port)

	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.HandleFunc("/socket", WebSocketHandler)

	mux.HandleFunc("/image", ImageHandler)

	// API routes
	mux.HandleFunc("/health", HealthHandler)

	// Static file server from embedded SPA build
	mux.Handle("/", NewEmbedHandler())

	srv := &http.Server{Addr: addr, Handler: loggingMiddleware(mux)}

	// Graceful shutdown
	go func() {
		log.Printf("Server listening on http://localhost%s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Wait for the interrupt signal
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &ResponseWriter{ResponseWriter: w, Status: 200}
		next.ServeHTTP(rw, r)
		dur := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.Status, dur)
	})
}
