package handlers

import (
	"io/fs"
	"net/http"

	"github.com/diegorezm/DBlockchain/internals/blockchain"
	"github.com/diegorezm/DBlockchain/internals/frontend"
	"github.com/diegorezm/DBlockchain/internals/frontend/pages/blocks_page"
	"github.com/diegorezm/DBlockchain/internals/frontend/pages/transactions_page"
	"github.com/diegorezm/DBlockchain/internals/frontend/pages/wallet_page"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FrontendHandler struct {
	blockchain *blockchain.Blockchain
}

func NewFrontendHandler(blockchain *blockchain.Blockchain) *FrontendHandler {
	return &FrontendHandler{blockchain}
}

func (h *FrontendHandler) GetIndexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/wallet", http.StatusMovedPermanently)
}

func (h *FrontendHandler) GetWalletPage(w http.ResponseWriter, r *http.Request) {
	var utxos []blockchain.UTXO

	publicKey := getPublicKeyFromCookies(r)

	if publicKey != "" {
		utxos = h.blockchain.GetUTXPoolByAddress(publicKey)
	}

	walletPage := wallet_page.WalletPage(publicKey, utxos)

	ctx := r.Context()
	if err := walletPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *FrontendHandler) GetCreateWalletPage(w http.ResponseWriter, r *http.Request) {
	walletPage := wallet_page.CreateWalletPage()

	ctx := r.Context()
	if err := walletPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *FrontendHandler) GetBlocksPage(w http.ResponseWriter, r *http.Request) {
	blocksPage := blocks_page.BlocksPage(h.blockchain.Chain)
	ctx := r.Context()
	if err := blocksPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *FrontendHandler) GetTransactionsPage(w http.ResponseWriter, r *http.Request) {
	publicKey := getPublicKeyFromCookies(r)

	transactionsPage := transactions_page.TransactionsPage(publicKey, h.blockchain.TransactionsMempool)

	ctx := r.Context()

	if err := transactionsPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *FrontendHandler) ServeAssets(r chi.Router) {
	r.Use(middleware.StripSlashes)
	if frontend.IsDev {
		// Serve from disk
		fileServer := http.FileServer(http.Dir("./internals/frontend/assets"))
		r.Handle("/*", http.StripPrefix("/assets", fileServer))
	} else {
		assetsFS, err := fs.Sub(frontend.EmbdedAssets, "assets")
		if err != nil {
			panic(err)
		}
		fileServer := http.FileServer(http.FS(assetsFS))
		r.Handle("/*", http.StripPrefix("/assets", fileServer))
	}
	fileServer := http.FileServer(http.Dir("./internals/frontend/assets"))
	r.Handle("/*", http.StripPrefix("/assets", fileServer))
}

func getPublicKeyFromCookies(r *http.Request) string {
	publicKeyCookie, err := r.Cookie("public-key")

	if err != nil {
		return ""
	} else {
		return publicKeyCookie.Value
	}

}
