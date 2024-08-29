package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request
		log.Printf("Received request from %s\n", r.RemoteAddr)
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)

		// Respond with a hello message
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Hello From Backend Server")
		log.Println("Replied with a hello message")
	})

	// Start the backend server
	port := ":8080" // Listen on port 8080
	log.Printf("Starting backend server on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
