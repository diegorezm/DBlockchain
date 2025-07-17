package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/frontend/components/alerts"
	"github.com/diegorezm/DBlockchain/internals/frontend/pages/blocks_page"
	"github.com/diegorezm/DBlockchain/internals/utils"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

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
	PrivateKey string  `schema:"private_key"`
	From       string  `schema:"from"`
	To         string  `schema:"to"`
	Amount     float64 `schema:"amount"`
}

func (bc *BlockchainClientHandler) AppendTransaction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError("Could not parse your request."), r.Context())
		return
	}

	var input appendTransactionInput
	if err := decoder.Decode(&input, r.PostForm); err != nil {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError("Could not parse your request."), r.Context())
		return
	}

	privKey, err := utils.DecodePrivateKey(input.PrivateKey)
	if err != nil {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError("Failed to parse private key"), r.Context())
		return
	}

	availableUTXOs := bc.blockchain.GetUTXPoolByAddress(input.From)

	var txIns []blockchain.TxIn
	var totalInput float64
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
		webutils.WriteTempl(w, http.StatusInternalServerError, alerts.AlertWarning("Insufficient funds."), r.Context())
		return
	}

	txOuts := []blockchain.TxOut{
		{Address: input.To, Amount: input.Amount},
	}
	if change := totalInput - input.Amount; change > 0 {
		txOuts = append(txOuts, blockchain.TxOut{Address: input.From, Amount: change})
	}

	txInput := blockchain.TransactionInput{
		TxIns:    txIns,
		TxOuts:   txOuts,
		IsSystem: false,
	}

	signedTx, err := blockchain.NewSignedTransaction(txInput, privKey)
	if err != nil {
		webutils.WriteTempl(w, http.StatusInternalServerError, alerts.AlertError("Signing failed."), r.Context())
		return
	}

	if err := bc.blockchain.AppendTransaction(signedTx); err != nil {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError(fmt.Sprintf("Failed to add transaction: %v", err)), r.Context())
		return
	}

	webutils.WriteTempl(w, http.StatusOK, alerts.AlertInfo("Transaction added to pool."), r.Context())
}

type buyCoinsRequest struct {
	PrivateKey string  `schema:"private_key"`
	To         string  `schema:"to"`
	Amount     float32 `schema:"amount"`
}

func (bc *BlockchainClientHandler) BuyCoins(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		webutils.WriteBadRequest(w, "Could not parse your request.")
		return
	}

	var input buyCoinsRequest
	err = decoder.Decode(&input, r.PostForm)

	if err != nil {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError(fmt.Sprintf("Failed to parse your request: %v", err)), r.Context())
		return
	}

	privKey, err := utils.DecodePrivateKey(input.PrivateKey)
	if err != nil {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError(fmt.Sprintf("Failed to parse your private key: %v", err)), r.Context())
		return
	}

	if input.Amount <= 0 {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError("The amount should be greater than 0."), r.Context())
		return
	}

	if input.To == "" {
		webutils.WriteTempl(w, http.StatusBadRequest, alerts.AlertError("Not a valid address."), r.Context())
		return
	}

	txInput := blockchain.TransactionInput{
		IsSystem: true,
		TxIns:    []blockchain.TxIn{},
		TxOuts: []blockchain.TxOut{
			{
				Address: input.To,
				Amount:  float64(input.Amount),
			},
		},
	}

	signedTx, err := blockchain.NewSignedTransaction(txInput, privKey)
	if err != nil {
		webutils.WriteTempl(w, http.StatusInternalServerError, alerts.AlertError(fmt.Sprintf("Failed sign your transaction: %v", err)), r.Context())
		return
	}

	if err := bc.blockchain.AppendTransaction(signedTx); err != nil {
		webutils.WriteTempl(w, http.StatusInternalServerError, alerts.AlertError(fmt.Sprintf("Failed add your transaction: %v", err)), r.Context())
		return
	}

	webutils.WriteTempl(w, http.StatusOK, alerts.AlertInfo("Your transaction was successfull! Now just wait for your another block to be mined."), r.Context())
}

func (bc *BlockchainClientHandler) ReplaceChain(w http.ResponseWriter, r *http.Request) {
	replaced, err := bc.blockchain.ReplaceChain()

	if err != nil {
		webutils.WriteTempl(w, http.StatusOK, alerts.AlertError(fmt.Sprintf("Failed to replace chain: %v", err)), r.Context())
		return
	}

	message := "Chain replacement attempted."

	if replaced {
		message = "Blockchain was successfully replaced."
	} else {
		message = "Blockchain was not replaced (current chain is valid and/or longer)."
	}

	webutils.WriteTempl(w, http.StatusOK, alerts.AlertInfo(message), r.Context())
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
	r.Post("/chain/replace", bc.ReplaceChain)
	r.Post("/chain/mine", bc.Mine)
	r.Post("/transactions/add", bc.AppendTransaction)
	r.Post("/transactions/buy", bc.BuyCoins)
}
