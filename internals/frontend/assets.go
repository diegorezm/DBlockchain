package frontend

import "embed"

//go:embed assets
var EmbdedAssets embed.FS
var IsDev = true
