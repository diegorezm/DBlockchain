package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/frontend"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

func main() {
	mux := http.NewServeMux()

	port := flag.Int("port", 3000, "Port to listen on (default 3000)")
	flag.Parse()

	if *port == 4040 {
		panic("This address is reserved for the server.")
	}

	addr := fmt.Sprintf(":%d", *port)
	fullAddr := fmt.Sprintf("http://localhost:%d", *port)
	serverAddr := "http://localhost:4040"

	err := registerNode(serverAddr, fullAddr)

	if err != nil {
		panic(err)
	}

	blockchain := bl.NewBlockchain(fullAddr)
	registerHandlers(mux, blockchain)

	loggedMux := webutils.LoggingMiddleware(mux)

	fmt.Printf("Client listening on port %s\n", addr)
	log.Fatalf("%v", http.ListenAndServe(addr, loggedMux))
}

func registerHandlers(mux *http.ServeMux, blockchain *bl.Blockchain) {
	blockchainHandler := bl.NewBlockchainHandler(blockchain)
	frontendHandler := frontend.NewFrontendHandler()

	// PAGES ROUTES
	mux.Handle("GET /", http.HandlerFunc(frontendHandler.GetIndexPage))
	mux.Handle("GET /assets/", http.HandlerFunc(frontendHandler.GetAssets))

	// API ROUTES
	mux.Handle("GET /api/chain", http.HandlerFunc(blockchainHandler.GetChain))
	mux.Handle("GET /api/chain/is_valid", http.HandlerFunc(blockchainHandler.IsValid))
	mux.Handle("GET /api/chain/replace", http.HandlerFunc(blockchainHandler.ReplaceChain))
	mux.Handle("POST /api/chain/mine", http.HandlerFunc(blockchainHandler.Mine))

	mux.Handle("POST /api/transactions/add", http.HandlerFunc(blockchainHandler.AddTransaction))
	mux.Handle("POST /api/transactions/add/bulk", http.HandlerFunc(blockchainHandler.AddTransactionBulk))

	mux.Handle("GET /ping", http.HandlerFunc(blockchainHandler.PingHandler))
}

func registerNode(serverAddr, nodeUrl string) error {
	reqUrl := fmt.Sprintf("%s/connect", serverAddr)

	body := fmt.Sprintf(`{"address": "%s"}`, nodeUrl)
	payload := []byte(body)

	res, err := http.Post(reqUrl, "application/json", bytes.NewBuffer(payload))

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status code: %d", res.StatusCode)
	}

	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}
