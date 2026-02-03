package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 2 * time.Second}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)
		log.Printf("[ServiceB] %s %s %d %s", r.Method, r.URL.Path, rec.status, time.Since(start))
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func callEchoHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing msg parameter"})
		return
	}

	resp, err := client.Get("http://localhost:8080/echo?msg=" + msg)
	if err != nil {
		log.Printf("[ServiceB] ERROR calling ServiceA: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "ServiceA unavailable"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ServiceB] ERROR reading ServiceA response: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to read ServiceA response"})
		return
	}

	var echoResp map[string]interface{}
	json.Unmarshal(body, &echoResp)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"service": "B",
		"from_a":  echoResp,
	})
}

func main() {
	http.HandleFunc("/health", loggingMiddleware(healthHandler))
	http.HandleFunc("/call-echo", loggingMiddleware(callEchoHandler))

	log.Println("[ServiceB] listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
