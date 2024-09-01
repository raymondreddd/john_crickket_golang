package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	addr := flag.String("addr", ":4000", "HTTP network address")

	log.Printf("Server is listening on %s", *addr)

	// Use ListenAndServe with a custom server to set timeouts
	server := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
