package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	addr := flag.String("addr", ":8989", "HTTP network address")
	flag.Parse()

	mux.HandleFunc("/", handle)

	log.Printf("Server is listening on %s", *addr)

	http.ListenAndServe(*addr, mux)
}
