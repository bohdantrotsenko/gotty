package server

import "embed"

// staticFS holds the web UI assets, embedded at compile time.
//
//go:embed all:static
var staticFS embed.FS
