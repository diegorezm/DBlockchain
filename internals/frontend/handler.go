package frontend

import (
	"net/http"

	"github.com/diegorezm/DBlockchain/internals/frontend/pages"
	webutils "github.com/diegorezm/DBlockchain/internals/web_utils"
)

// //go:embed assets
// var embdedAssets embed.FS

type frontendHandler struct {
}

func NewFrontendHandler() *frontendHandler {
	return &frontendHandler{}
}

func (h *frontendHandler) GetIndexPage(w http.ResponseWriter, r *http.Request) {
	currentPath := r.URL.Path
	index := pages.Index(currentPath)
	ctx := r.Context()
	if err := index.Render(ctx, w); err != nil {
		webutils.WriteInternalServerError(w, err.Error())
	}
}

func (h *frontendHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/assets/", http.FileServer(http.Dir("./internals/frontend/assets"))).ServeHTTP(w, r)
}
