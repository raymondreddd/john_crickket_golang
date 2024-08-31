package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	log.Printf("Server is listening on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
