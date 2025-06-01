package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
)

func main() {
	mux := http.NewServeMux()
	blockchain := bl.NewBlockchain()

	port := flag.Int("port", 4000, "Port to listen on")
	flag.Parse()

	blockchainHandlers(mux, blockchain)

	addr := fmt.Sprintf(":%d", *port)
	fullAddr := fmt.Sprintf("http://127.0.0.1%s", addr)
	blockchain.AppendNode(fullAddr)

	fmt.Printf("Server listening on port %s\n", addr)
	log.Fatalf("%v", http.ListenAndServe(addr, mux))
}

func blockchainHandlers(mux *http.ServeMux, blockchain *bl.Blockchain) {
	blockchainHandler := bl.NewHandler(blockchain)

	mux.Handle("GET /chain", http.HandlerFunc(blockchainHandler.GetChain))
	mux.Handle("GET /chain/is_valid", http.HandlerFunc(blockchainHandler.IsValid))
	mux.Handle("POST /chain/mine", http.HandlerFunc(blockchainHandler.Mine))
	mux.Handle("POST /chain/replace", http.HandlerFunc(blockchainHandler.ReplaceChain))

	mux.Handle("POST /transactions/add", http.HandlerFunc(blockchainHandler.AddTransaction))
	mux.Handle("POST /transactions/add/bulk", http.HandlerFunc(blockchainHandler.AddTransactionBulk))

	mux.Handle("POST /nodes/add", http.HandlerFunc(blockchainHandler.AddNode))
	mux.Handle("POST /nodes/add/bulk", http.HandlerFunc(blockchainHandler.AddNodeBulk))
}
