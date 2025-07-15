package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/frontend/pages/blocks_page"
	"github.com/diegorezm/DBlockchain/internals/utils"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
	"github.com/go-chi/chi/v5"
)

type BlockchainClientHandler struct {
	blockchain *blockchain.Blockchain
}

func NewBlockchainClientHandler(bl *blockchain.Blockchain) *BlockchainClientHandler {
	return &BlockchainClientHandler{
		blockchain: bl,
	}
}

func (bc *BlockchainClientHandler) GetChain(w http.ResponseWriter, r *http.Request) {
	chain := bc.blockchain.Chain
	webutils.WriteJSON(w, 200, chain, "Blocks fetched")
}

func (bc *BlockchainClientHandler) Mine(w http.ResponseWriter, r *http.Request) {
	if err := bc.blockchain.AppendBlock(); err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Failed to mine new block: %v", err))
		return
	}
	time.Sleep(time.Duration(5000))
	chain := bc.blockchain.Chain
	table := blocks_page.BlocksTable(chain)
	w.Header().Set("Content-Type", "text/html")
	table.Render(r.Context(), w)
	// webutils.WriteJSON[any](w, http.StatusCreated, nil, "New block mined successfully!")
}

type appendTransactionInput struct {
	PrivateKey string `json:"private_key"`
	blockchain.TransactionInput
}

func (bc *BlockchainClientHandler) AppendTransaction(w http.ResponseWriter, r *http.Request) {
	input, err := webutils.ParseJSON[appendTransactionInput](r.Body)
	if err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	privKey, err := utils.DecodePrivateKey(input.PrivateKey)
	if err != nil {
		webutils.WriteBadRequest(w, "Failed to parse private key")
		return
	}

	tx, err := blockchain.NewSignedTransaction(input.TransactionInput, privKey)

	if err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Something went wrong: %v", err))
		return
	}

	err = bc.blockchain.AppendTransaction(tx)

	if err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid transaction: %v", err))
		return
	}
}

func (bc *BlockchainClientHandler) ReplaceChain(w http.ResponseWriter, r *http.Request) {
	replaced, err := bc.blockchain.ReplaceChain()
	if err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Failed to replace chain: %v", err))
		return
	}

	respData := map[string]bool{
		"replaced": replaced,
	}

	message := "Chain replacement attempted."
	if replaced {
		message = "Blockchain was successfully replaced."
	} else {
		message = "Blockchain was not replaced (current chain is valid and/or longer)."
	}

	webutils.WriteSuccess(w, respData, message)
}

func (bc *BlockchainClientHandler) IsChainValid(w http.ResponseWriter, r *http.Request) {
	isValid := blockchain.IsChainValid(bc.blockchain.Chain)

	respData := map[string]bool{
		"isValid": isValid,
	}
	webutils.WriteSuccess(w, respData, "Blockchain validation status.")
}

func (bc *BlockchainClientHandler) Register(r chi.Router) {
	r.Get("/chain", bc.GetChain)
	r.Get("/chain/is_valid", bc.IsChainValid)
	r.Get("/chain/replace", bc.ReplaceChain)
	r.Post("/chain/mine", bc.Mine)
	r.Post("/transactions/add", bc.AppendTransaction)
}
