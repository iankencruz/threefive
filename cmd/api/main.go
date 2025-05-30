package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {

	r := chi.NewRouter()

	log.Println("🚀 Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("❌ Server error:", err)
	}

}
