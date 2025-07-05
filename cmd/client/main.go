package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/frontend"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

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
	registerHandlers(r, blockchain)

	fmt.Printf("Client listening on port %s\n", addr)
	log.Fatalf("%v", http.ListenAndServe(addr, r))
}

func registerHandlers(r *chi.Mux, blockchain *bl.Blockchain) {
	blockchainHandler := bl.NewBlockchainHandler(blockchain)
	frontendHandler := frontend.NewFrontendHandler()

	// PAGES
	r.Route("/", func(r chi.Router) {
		r.Get("/", frontendHandler.GetIndexPage)
		r.Get("/wallet", frontendHandler.GetWalletPage)
		r.Get("/blocks", frontendHandler.GetBlocksPage)
		r.Get("/transactions", frontendHandler.GetTransactionsPage)
	})

	r.Route("/assets", func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		fileServer := http.FileServer(http.Dir("./internals/frontend/assets"))
		r.Handle("/*", http.StripPrefix("/assets", fileServer))
	})

	// API
	r.Route("/api", func(r chi.Router) {
		r.Get("/chain", blockchainHandler.GetChain)
		r.Get("/chain/is_valid", blockchainHandler.IsValid)
		r.Get("/chain/replace", blockchainHandler.ReplaceChain)
		r.Post("/chain/mine", blockchainHandler.Mine)

		r.Post("/transactions/add", blockchainHandler.AddTransaction)
		r.Post("/transactions/add/bulk", blockchainHandler.AddTransactionBulk)
	})

	// Other
	r.Get("/ping", blockchainHandler.PingHandler)

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
