package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"

	bl "github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/handlers"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
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
	blockchainHandler := handlers.NewBlockchainClientHandler(blockchain)
	walletHandler := handlers.NewWalletHandler(blockchain)
	frontendHandler := handlers.NewFrontendHandler()

	// PAGES
	r.Route("/", func(r chi.Router) {
		r.Get("/", frontendHandler.GetIndexPage)
		r.Get("/wallet", frontendHandler.GetWalletPage)
		r.Get("/blocks", frontendHandler.GetBlocksPage)
		r.Get("/transactions", frontendHandler.GetTransactionsPage)
	})

	r.Route("/assets", frontendHandler.ServeAssets)

	// API
	r.Route("/api", func(r chi.Router) {
		blockchainHandler.Register(r)
		walletHandler.Register(r)
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		webutils.WriteSuccess[any](w, nil, "Pong!")
	})

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
