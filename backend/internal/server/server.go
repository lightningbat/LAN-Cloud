package server

import (
	"lan-cloud/internal/server/httpserver"
	"log"
	"net/http"
)

func Start() {
	mux := http.NewServeMux()

	httpserver.RegisterRoutes(mux)

	handleWithCORS := httpserver.CORSMiddleware(mux)

	port := ":8080"
	log.Println("Server running at http://localhost" + port)
	err := http.ListenAndServe(port, handleWithCORS)
	if err != nil {
		log.Fatal(err)
	}
}