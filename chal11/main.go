package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.URL.Path
	proto := r.Proto

	// for server log
	log.Printf("Received request: Method = %s, Path = %s, HTTP Version = %s", method, path, proto)

	// client res log
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nRequested path: %s\r\n", path)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, response)
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	log.Printf("Server is listening on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
