package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting AI Interview Platform API Server...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
	http.ListenAndServe(":8080", nil)
}
