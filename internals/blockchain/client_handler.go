package blockchain

import (
	"encoding/json"
	"fmt"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
	"net/http"
)

type BlockchainHandler struct {
	blockchain *Blockchain
}

func NewBlockchainHandler(blockchain *Blockchain) *BlockchainHandler {
	return &BlockchainHandler{blockchain: blockchain}
}

func (h *BlockchainHandler) GetChain(w http.ResponseWriter, r *http.Request) {
	chain := h.blockchain.GetChain()
	webutils.WriteSuccess(w, chain, "Blockchain retrieved successfully")
}

func (h *BlockchainHandler) Mine(w http.ResponseWriter, r *http.Request) {
	if err := h.blockchain.AppendBlock(); err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Failed to mine new block: %v", err))
		return
	}
	webutils.WriteJSON[any](w, http.StatusCreated, nil, "New block mined successfully!")
}

func (h *BlockchainHandler) IsValid(w http.ResponseWriter, r *http.Request) {
	valid := isChainValid(h.blockchain.chain)

	respData := map[string]bool{
		"valid": valid,
	}
	webutils.WriteSuccess(w, respData, "Blockchain validation status.")
}

func (h *BlockchainHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionInsert TransactionInsert
	if err := json.NewDecoder(r.Body).Decode(&transactionInsert); err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}
	h.blockchain.AppendTransaction(transactionInsert)
	webutils.WriteJSON[any](w, http.StatusCreated, nil, "Transaction added successfully.")
}

func (h *BlockchainHandler) AddTransactionBulk(w http.ResponseWriter, r *http.Request) {
	var transactionBulkRequest TransactionBulkRequest

	if err := json.NewDecoder(r.Body).Decode(&transactionBulkRequest); err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload for bulk transactions: %v", err))
		return
	}

	for _, t := range transactionBulkRequest.Transactions {
		h.blockchain.AppendTransaction(t)
	}

	webutils.WriteJSON[any](w, http.StatusCreated, nil, "Transactions added successfully.")
}

func (h *BlockchainHandler) ReplaceChain(w http.ResponseWriter, r *http.Request) {
	replaced, err := h.blockchain.replaceChain()
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

func (s *BlockchainHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	webutils.WriteSuccess[any](w, nil, "Pong!")
}
