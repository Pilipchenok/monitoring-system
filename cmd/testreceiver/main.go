package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		fmt.Printf("Received: %s\n", string(body))
		w.WriteHeader(http.StatusOK)
	})

	log.Println("Test receiver on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
