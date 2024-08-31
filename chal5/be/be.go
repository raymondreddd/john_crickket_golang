package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request from", r.RemoteAddr)
	fmt.Println(r.Method, r.URL, r.Proto)
	fmt.Println("Host:", r.Host)
	fmt.Println("User-Agent:", r.UserAgent())
	fmt.Println("Accept:", r.Header.Get("Accept"))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello From Backend Server")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting backend server on port 8080")
	http.ListenAndServe(":8080", nil)
}
