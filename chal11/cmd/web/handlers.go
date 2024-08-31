package main

import (
	"fmt"
	"html/template" // New import
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	method := r.Method
	path := r.URL.Path
	proto := r.Proto

	// for server log
	log.Printf("Received request: Method = %s, Path = %s, HTTP Version = %s", method, path, proto)

	// client res log
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n\r\nRequested path: %s\r\n", path)

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, response)
	// Use the template.ParseFiles() function to read the template file into a
	// template set. If there's an error, we log the detailed error message and
	// the http.Error() function to send a generic 500 Internal Server Error
	// response to the user.
	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	// We then use the Execute() method on the template set to write the templa
	// content as the response body. The last parameter to Execute() represents
	// dynamic data that we want to pass in, which for now we'll leave as nil.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
