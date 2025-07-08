package handlers

import (
	"log"
	"net/http"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/utils"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
	"github.com/go-chi/chi/v5"
)

type WalletHandler struct {
	blockchain *blockchain.Blockchain
}

func NewWalletHandler(blockchain *blockchain.Blockchain) *WalletHandler {
	return &WalletHandler{blockchain: blockchain}
}

func (wh *WalletHandler) GetUTXOsByAddress(w http.ResponseWriter, r *http.Request) {
	address := chi.URLParam(r, "address")
	utxos := wh.blockchain.GetUTXPoolByAddress(address)
	webutils.WriteJSON(w, 200, utxos, "Here are the unspent transactions for the given address.")
}

func (wh *WalletHandler) Generate(w http.ResponseWriter, r *http.Request) {
	priv, err := utils.GenerateKeyPair()

	if err != nil {
		log.Printf("%v", err)
		webutils.WriteInternalServerError(w, "Something went wrong while generating your key pair")
		return
	}

	keypair, err := utils.EncodeKeyPair(priv)
	if err != nil {
		log.Printf("%v", err)
		webutils.WriteInternalServerError(w, "Failed to encode key")
		return
	}

	webutils.WriteSuccess(w, keypair, "Key pair generated")
}

func (wh *WalletHandler) Register(r chi.Router) {
	r.Get("/wallet/generate", wh.Generate)
	r.Get("/wallet/utxos/{address}", wh.GetUTXOsByAddress)
}
