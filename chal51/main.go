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

	err := loadBannedHosts("forbidden-hosts.txt")
	if err != nil {
		log.Fatalf("Error loading banned hosts: %v", err)
	}

	mux.HandleFunc("/", handleRequestAndRedirect)

	log.Printf("Server is listening on %s", *addr)

	http.ListenAndServe(*addr, mux)
}
