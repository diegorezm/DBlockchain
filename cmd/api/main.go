package main

import (
	"fmt"
	"net/http"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
)

func main() {
	mux := &http.ServeMux{}

	handlers(mux)

	addr := ":4000"
	fmt.Printf("Server listening on http://localhost%s/\n", addr)
	http.ListenAndServe(addr, mux)
}

func handlers(mux *http.ServeMux) {
	blockchain := bl.NewBlockchain()
	blockchainHandler := bl.NewHandler(blockchain)

	mux.Handle("GET /get_chain", http.HandlerFunc(blockchainHandler.GetChain))
	mux.Handle("POST /mine", http.HandlerFunc(blockchainHandler.Mine))
}
