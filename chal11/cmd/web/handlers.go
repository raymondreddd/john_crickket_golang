package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Log request details and simulate a long processing time
	method := r.Method
	path := r.URL.Path
	proto := r.Proto
	threadID := time.Now().UnixNano()

	// Log the thread ID (goroutine ID)
	log.Printf("Thread Id: %d, Received request: Method = %s, Path = %s, HTTP Version = %s", threadID, method, path, proto)

	// Simulate long processing time (5 seconds)
	time.Sleep(5 * time.Second)

	// Prepare and send the response
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "HTTP/1.1 200 OK\r\n\r\nRequested path: %s\r\n", path)

	ts, err := template.ParseFiles("../../ui/html/home.page.tmpl")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
