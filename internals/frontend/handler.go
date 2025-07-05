package frontend

import (
	"embed"
	"net/http"

	"github.com/diegorezm/DBlockchain/internals/frontend/pages"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

// //go:embed assets
var EmbdedAssets embed.FS

type frontendHandler struct {
}

func NewFrontendHandler() *frontendHandler {
	return &frontendHandler{}
}

func (h *frontendHandler) GetIndexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/wallet", http.StatusMovedPermanently)
}

func (h *frontendHandler) GetWalletPage(w http.ResponseWriter, r *http.Request) {
	walletPage := pages.WalletPage()
	ctx := r.Context()
	if err := walletPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *frontendHandler) GetBlocksPage(w http.ResponseWriter, r *http.Request) {
	blocksPage := pages.BlocksPage()
	ctx := r.Context()
	if err := blocksPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *frontendHandler) GetTransactionsPage(w http.ResponseWriter, r *http.Request) {
	transactionsPage := pages.TransactionsPage()
	ctx := r.Context()
	if err := transactionsPage.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}
