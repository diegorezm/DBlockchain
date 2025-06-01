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
	if err := h.blockchain.AppendBlock(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) IsValid(w http.ResponseWriter, r *http.Request) {
	valid := isChainValid(h.blockchain.chain)
	resp := map[string]bool{
		"valid": valid,
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *Handler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionInsert TransactionInsert
	if err := json.NewDecoder(r.Body).Decode(&transactionInsert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.blockchain.AppendTransaction(transactionInsert)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) AddTransactionBulk(w http.ResponseWriter, r *http.Request) {
	var transactionBulkRequest TransactionBulkRequest

	if err := json.NewDecoder(r.Body).Decode(&transactionBulkRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, t := range transactionBulkRequest.Transactions {
		h.blockchain.AppendTransaction(t)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) AddNode(w http.ResponseWriter, r *http.Request) {
	var nodeInsert struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&nodeInsert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.blockchain.AppendNode(nodeInsert.Address)

	resp := map[string][]Node{
		"nodes": h.blockchain.nodes.ToSlice(),
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *Handler) AddNodeBulk(w http.ResponseWriter, r *http.Request) {
	var nodeBulkInsert NodeBulkRequest

	if err := json.NewDecoder(r.Body).Decode(&nodeBulkInsert); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for _, addr := range nodeBulkInsert.Nodes {
		h.blockchain.AppendNode(addr)
	}

	resp := map[string][]Node{
		"nodes": h.blockchain.nodes.ToSlice(),
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (h *Handler) ReplaceChain(w http.ResponseWriter, r *http.Request) {
	replaced := h.blockchain.replaceChain()

	resp := map[string]bool{
		"replaced": replaced,
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
