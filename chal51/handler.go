package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

// Handle each request, checking if the host is banned
func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if _, banned := BannedHosts[host]; banned {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		msg := fmt.Sprintf("Website not allowed: %s\n", host)
		w.Write([]byte(msg))
		log.Printf("Blocked request to banned host: %s", host)
		return
	}

	// Forward request to the target if not banned
	proxyHanlder(w, r)
}

func proxyHanlder(w http.ResponseWriter, r *http.Request) {
	// Extract client IP for X-Forwarded-For header
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Unable to determine client IP", http.StatusInternalServerError)
		return
	}

	/*
		If we don't use the httputil.ReverseProxy then we can use below

		Parse the incoming request.
		Recreate the request for the target server, ensuring headers are copied and adjusted.
		Forward the request to the target server using http.Client.
		Read the server’s response, and copy the response back to the client, adjusting headers as needed.
		targetURL := fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)

		client := &http.Client{}
		req, err := http.NewRequest(r.Method, targetURL, r.Body)
	*/

	// Remove hop-by-hop headers (which apply to a single connection)
	hopByHopHeaders := []string{
		"Proxy-Connection",
		"Connection",
		"Keep-Alive",
		"Proxy-Authenticate",
		"Proxy-Authorization",
		"TE",
		"Trailers",
		"Transfer-Encoding",
		"Upgrade"}
	// hopByHopHeaders := []string{
	//     "Proxy-Connection", "Connection", "Keep-Alive", "Proxy-Authenticate",
	//     "Proxy-Authorization", "TE", "Trailers", "Transfer-Encoding", "Upgrade",
	//     "X-Forwarded-For", "X-Forwarded-Host", "X-Forwarded-Proto", "Forwarded", "Via",
	// }
	for _, h := range hopByHopHeaders {
		r.Header.Del(h)
	}

	// Add X-Forwarded-For header
	r.Header.Set("X-Forwarded-For", clientIP)

	// Modify the request for proxy use (ensure full URL is used)
	r.URL.Scheme = "http"

	// if conn made over Transport Security Layer
	if r.TLS != nil {
		r.URL.Scheme = "https"
	}
	r.URL.Host = r.Host

	// Forward request
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = r.URL.Scheme
			req.URL.Host = r.Host
			req.Host = r.Host
		},
	}
	proxy.ServeHTTP(w, r)
}
