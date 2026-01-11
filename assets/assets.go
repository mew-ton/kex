package assets

import "embed"

//go:embed all:templates manifest.json
var Assets embed.FS
