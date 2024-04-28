package main

import (
	"git.garena.com/shopee/insurance/insurance-backend/insurance-hub/moving-window/src/server"
	_ "git.garena.com/shopee/insurance/insurance-backend/insurance-hub/moving-window/src/server/impl"
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
