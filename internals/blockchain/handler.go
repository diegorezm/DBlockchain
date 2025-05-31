package blockchain

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	blockchain *Blockchain
}

func NewHandler(blockchain *Blockchain) *Handler {
	return &Handler{blockchain: blockchain}
}

func (h *Handler) GetChain(w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := json.Marshal(h.blockchain.chain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *Handler) Mine(w http.ResponseWriter, r *http.Request) {
	transactionInsert := TransactionInsert{
		From:   "Guy",
		To:     "Another guy",
		Amount: 1,
	}
	h.blockchain.AppendTransaction(transactionInsert)
	h.blockchain.AppendTransaction(transactionInsert)
	if err := h.blockchain.AppendBlock(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
