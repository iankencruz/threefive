package main

import (
	"log"
	"net/http"
	"time"

	"github.com/iankencruz/threefive/internal/application"
	"github.com/iankencruz/threefive/internal/routes"
)

func main() {

	app := application.NewApp()

	r := routes.Routes(app)

	log.Println("🚀 Server running at http://localhost:8080")
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("❌ Server error:", err)
	}

}
