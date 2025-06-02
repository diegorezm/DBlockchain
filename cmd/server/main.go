package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

func main() {
	mux := http.NewServeMux()
	handler := bl.NewServerHandler()

	registerEndpoints(mux, handler)

	port := 4040
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server listening on port %s\n", addr)

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			handler.PingNodes()
		}
	}()

	loggedMux := webutils.LoggingMiddleware(mux)
	log.Fatalf("%v", http.ListenAndServe(addr, loggedMux))
}

func registerEndpoints(mux *http.ServeMux, handler *bl.ServerHandler) {
	mux.Handle("GET /nodes", http.HandlerFunc(handler.GetNodes))

	mux.Handle("POST /connect", http.HandlerFunc(handler.ConnectNode))
	mux.Handle("POST /disconnect", http.HandlerFunc(handler.DisconnectNode))

	mux.Handle("POST /ping", http.HandlerFunc(handler.PingHandler))
}
