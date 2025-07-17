package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/frontend/pages/wallet_page"
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

	page := wallet_page.PublicAndPrivateKeyGeneration(keypair.PublicKey, keypair.PrivateKey, true)
	w.Header().Set("Content-Type", "text/html")
	page.Render(r.Context(), w)
}

func (wh *WalletHandler) SavePubKey(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		webutils.WriteBadRequest(w, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	v := r.Form.Get("pubKey")

	cookie := http.Cookie{
		Name:     "public-key",
		Value:    v,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/wallet", http.StatusSeeOther)
}

func (wh *WalletHandler) ForgetPublicKey(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:     "public-key",
		Value:    "deleted",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/wallet", http.StatusSeeOther)
}

func (wh *WalletHandler) Register(r chi.Router) {
	r.Post("/wallet/generate", wh.Generate)
	r.Post("/wallet/save-key", wh.SavePubKey)
	r.Post("/wallet/forget-key", wh.ForgetPublicKey)
	r.Get("/wallet/utxos/{address}", wh.GetUTXOsByAddress)
}
