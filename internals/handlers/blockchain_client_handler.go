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
}

type appendTransactionInput struct {
	PrivateKey string  `json:"private_key"`
	From       string  `json:"from"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
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

	availableUTXOs := bc.blockchain.GetUTXPoolByAddress(input.From)

	var txIns []blockchain.TxIn
	var totalInput float64
	// For each availiable UTXO, we append a TxIn. We do this until there is enougth money
	for _, utxo := range availableUTXOs {
		txIns = append(txIns, blockchain.TxIn{
			TxOutId:    utxo.TxId,
			TxOutIndex: utxo.Index,
		})

		totalInput += utxo.Output.Amount

		if totalInput >= input.Amount {
			break
		}
	}

	if totalInput < input.Amount {
		webutils.WriteBadRequest(w, "Insufficient funds")
		return
	}

	var txOuts []blockchain.TxOut

	// Primary recipient
	txOuts = append(txOuts, blockchain.TxOut{
		Address: input.To,
		Amount:  input.Amount,
	})

	change := totalInput - input.Amount
	if change > 0 {
		txOuts = append(txOuts, blockchain.TxOut{
			Address: input.From,
			Amount:  change,
		})
	}

	txInput := blockchain.TransactionInput{
		TxIns:    txIns,
		TxOuts:   txOuts,
		IsSystem: false,
	}

	signedTx, err := blockchain.NewSignedTransaction(txInput, privKey)
	if err != nil {
		webutils.WriteInternalServerError(w, err.Error())
		return
	}

	if err := bc.blockchain.AppendTransaction(signedTx); err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Failed to add transaction: %v", err))
		return
	}
	webutils.WriteSuccess(w, signedTx, "Transaction added to pool")
}

type buyCoinsRequest struct {
	PrivateKey string  `json:"private_key"`
	To         string  `json:"to"`
	Amount     float64 `json:"amount"`
}

func (bc *BlockchainClientHandler) BuyCoins(w http.ResponseWriter, r *http.Request) {
	input, err := webutils.ParseJSON[buyCoinsRequest](r.Body)
	if err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	privKey, err := utils.DecodePrivateKey(input.PrivateKey)
	if err != nil {
		webutils.WriteBadRequest(w, "Failed to parse private key")
		return
	}

	if input.Amount <= 0 {
		webutils.WriteBadRequest(w, "Amount must be greater than zero")
		return
	}

	if input.To == "" {
		webutils.WriteBadRequest(w, "Recipient public key is required")
		return
	}

	txInput := blockchain.TransactionInput{
		IsSystem: true,
		TxIns:    []blockchain.TxIn{},
		TxOuts: []blockchain.TxOut{
			{
				Address: input.To,
				Amount:  input.Amount,
			},
		},
	}

	signedTx, err := blockchain.NewSignedTransaction(txInput, privKey)
	if err != nil {
		webutils.WriteInternalServerError(w, err.Error())
		return
	}

	if err := bc.blockchain.AppendTransaction(signedTx); err != nil {
		webutils.WriteInternalServerError(w, fmt.Sprintf("Failed to add transaction: %v", err))
		return
	}
	webutils.WriteSuccess(w, signedTx, "Transaction added to pool")
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
	r.Post("/transactions/buy", bc.BuyCoins)
}
