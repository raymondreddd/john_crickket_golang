package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf("Requested path: %s", r.URL.Path)

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
