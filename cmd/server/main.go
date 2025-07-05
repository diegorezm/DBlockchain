package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
)

func main() {
	r := chi.NewRouter()
	handler := bl.NewServerHandler()

	registerEndpoints(r, handler)

	port := 4040
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on port %s\n", addr)

	// Periodic ping
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			handler.PingNodes()
		}
	}()

	log.Fatalf("%v", http.ListenAndServe(addr, r))
}

func registerEndpoints(r chi.Router, handler *bl.ServerHandler) {
	r.Get("/nodes", handler.GetNodes)

	r.Post("/connect", handler.ConnectNode)
	r.Post("/disconnect", handler.DisconnectNode)

	r.Post("/ping", handler.PingHandler)
}
