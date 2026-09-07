package main

import (
	"log"
	"net/http"
	"os"

	"github.com/zdub0is/adhd-productivity-app/api/internal/health"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
