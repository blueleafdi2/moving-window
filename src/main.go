package main

import (
	"github.com/blueleafdi2/moving-window/src/server"
	_ "github.com/blueleafdi2/moving-window/src/server/impl"
	"log"
	"net/http"
)


func main() {
	srv := server.RefService()
	http.HandleFunc("/api/requests", srv.CountHandler)

	log.Println("Server is running on http://localhost:8081/api/requests")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("HTTP server error: %v\n", err)
	}
}
