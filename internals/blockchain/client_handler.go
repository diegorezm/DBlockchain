package blockchain

import (
	"encoding/json"
	"fmt"
	"net/http"

	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

type ClientHandler struct {
	blockchain *Blockchain
}

// NewClientHandler creates and returns a new ClientHandler.
func NewClientHandler(blockchain *Blockchain) *ClientHandler {
	return &ClientHandler{blockchain: blockchain}
}

// GetChain handles requests to retrieve the entire blockchain.
func (h *ClientHandler) GetChain(w http.ResponseWriter, r *http.Request) {
	chain := h.blockchain.GetChain() // Assuming GetChain() method exists on Blockchain
	webutils.WriteSuccess(w, chain, "Blockchain retrieved successfully")
}

// Mine handles requests to mine a new block.
func (h *ClientHandler) Mine(w http.ResponseWriter, r *http.Request) {
	if err := h.blockchain.AppendBlock(); err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Failed to mine new block: %v", err))
		return
	}
	webutils.WriteJSON[any](w, http.StatusCreated, nil, "New block mined successfully!")
}

// IsValid checks if the blockchain is valid.
func (h *ClientHandler) IsValid(w http.ResponseWriter, r *http.Request) {
	valid := isChainValid(h.blockchain.chain)

	respData := map[string]bool{
		"valid": valid,
	}
	webutils.WriteSuccess(w, respData, "Blockchain validation status.")
}

// AddTransaction handles adding a single new transaction.
func (h *ClientHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionInsert TransactionInsert
	if err := json.NewDecoder(r.Body).Decode(&transactionInsert); err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}
	h.blockchain.AppendTransaction(transactionInsert)
	webutils.WriteJSON[any](w, http.StatusCreated, nil, "Transaction added successfully.")
}

// AddTransactionBulk handles adding multiple transactions in a single request.
func (h *ClientHandler) AddTransactionBulk(w http.ResponseWriter, r *http.Request) {
	var transactionBulkRequest TransactionBulkRequest

	if err := json.NewDecoder(r.Body).Decode(&transactionBulkRequest); err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload for bulk transactions: %v", err))
		return
	}

	for _, t := range transactionBulkRequest.Transactions {
		h.blockchain.AppendTransaction(t) // Assuming AppendTransaction can handle TransactionInsert
	}

	webutils.WriteJSON[any](w, http.StatusCreated, nil, "Transactions added successfully.")
}

// ReplaceChain handles replacing the current blockchain with a longer, valid one from peers.
func (h *ClientHandler) ReplaceChain(w http.ResponseWriter, r *http.Request) {
	replaced, err := h.blockchain.replaceChain() // Assuming ReplaceChain() method exists on Blockchain
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

func (s *ClientHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	webutils.WriteSuccess[any](w, nil, "Pong!")
}
