package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the hand
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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

func startBackendServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// Use the Header().Set() method to add an 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second parameter is the header value.
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed) // same as above
		return
	}
	w.Write([]byte("Started backend"))
}
func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/start", startBackendServer)
	// Use the http.ListenAndServe() function to start a new web server. We pas
	// two parameters: the TCP network address to listen on (in this case ":4000
	// and the servemux we just created. If http.ListenAndServe() returns an er
	// we use the log.Fatal() function to log the error message and exit.
	log.Println("Starting server on :80")

	if err := http.ListenAndServe("localhost:80", mux); err != nil {
		fmt.Println(err.Error())
		log.Fatal(err)
	}

	// IP from request, host, user-agent, accept
	fmt.Print("Recieved request from")
}
