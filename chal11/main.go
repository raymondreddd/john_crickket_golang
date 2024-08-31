package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Construct the response
	response := fmt.Sprintf("Requested path: %s", r.URL.Path)

	// Send the response
	fmt.Fprintln(w, response)
}

func main() {
	// Register the handler function for all paths
	http.HandleFunc("/", handler)

	// Start the server on port 80
	fmt.Println("Server is listening on 127.0.0.1:80")
	err := http.ListenAndServe("127.0.0.1:80", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
