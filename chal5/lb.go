package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

var (
	backendServers = []string{
		"http://localhost:8080",
		"http://localhost:8081",
	}
	currentServer uint32
)

func getNextServer() string {
	index := atomic.AddUint32(&currentServer, 1)
	return backendServers[(int(index)-1)%len(backendServers)]
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request from", r.RemoteAddr)
	fmt.Println(r.Method, r.URL, r.Proto)
	fmt.Println("Host:", r.Host)
	fmt.Println("User-Agent:", r.UserAgent())
	fmt.Println("Accept:", r.Header.Get("Accept"))

	backendServer := getNextServer()

	proxyReq, err := http.NewRequest(r.Method, backendServer+r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	proxyReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting load balancer on port 8000")
	http.ListenAndServe(":8000", nil)
}
