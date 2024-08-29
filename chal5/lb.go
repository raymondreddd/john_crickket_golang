package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func handleRoute(path string) *bytes.Buffer {
	buf := bytes.Buffer{}

	switch path {
	case "/health":
		buf.Write([]byte("HTTP/1.1 204 No Content\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("\r\n"))
	case "/":
		buf.Write([]byte("HTTP/1.1 200 OK\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("Content-Length: 27\r\n"))
		buf.Write([]byte("\r\n"))
		buf.Write([]byte("Hello From Backend Server\r\n"))
	default:
		buf.Write([]byte("HTTP/1.1 404 Not Found\r\n"))
		buf.Write([]byte("Connection: close\r\n"))
		buf.Write([]byte("\r\n"))
	}

	return &buf
}

func startBackendServer(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the response header map
		w.Header().Set("Allow", "GET")
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // same as above 2 comments
		return
	}
	w.Write([]byte("Hello from Snippetbox"))

	backendURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to parse backend URL: %v", err)
	}

	// Create a reverse proxy to forward requests to the backend server
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	log.Printf("Received request from %s\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)

	// Forward the request to the backend server
	proxy.ServeHTTP(w, r)
}

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", startBackendServer)
	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an er
	// we use the log.Fatal() function to log the error message and exit.
	port := ":80" // Listen on port 80
	log.Printf("Starting load balancer on port %s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// IP from request, host, user-agent, accept
	fmt.Print("Recieved request from")
}
