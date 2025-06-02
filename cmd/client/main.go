package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"[%s] %s %s from %s - User-Agent: %s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.Header.Get("User-Agent"),
		)
		next.ServeHTTP(w, r)
		log.Printf(
			"Request for %s %s completed in %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func main() {
	mux := http.NewServeMux()
	blockchain := bl.NewBlockchain()

	port := flag.Int("port", 4000, "Port to listen on")
	flag.Parse()

	blockchainHandlers(mux, blockchain)

	addr := fmt.Sprintf(":%d", *port)

	loggedMux := loggingMiddleware(mux)

	fmt.Printf("Server listening on port %s\n", addr)
	log.Fatalf("%v", http.ListenAndServe(addr, loggedMux))
}

func blockchainHandlers(mux *http.ServeMux, blockchain *bl.Blockchain) {
	blockchainHandler := bl.NewHandler(blockchain)

	mux.Handle("GET /chain", http.HandlerFunc(blockchainHandler.GetChain))
	mux.Handle("GET /chain/is_valid", http.HandlerFunc(blockchainHandler.IsValid))
	mux.Handle("GET /chain/replace", http.HandlerFunc(blockchainHandler.ReplaceChain))
	mux.Handle("POST /chain/mine", http.HandlerFunc(blockchainHandler.Mine))

	mux.Handle("POST /transactions/add", http.HandlerFunc(blockchainHandler.AddTransaction))
	mux.Handle("POST /transactions/add/bulk", http.HandlerFunc(blockchainHandler.AddTransactionBulk))

	mux.Handle("POST /nodes/add", http.HandlerFunc(blockchainHandler.AddNode))
	mux.Handle("POST /nodes/add/bulk", http.HandlerFunc(blockchainHandler.AddNodeBulk))
}
