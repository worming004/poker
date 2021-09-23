package main

import (
	"embed"
	_ "embed"
	"net/http"
)

//go:embed static
var fs embed.FS

func getStaticHandler() http.Handler {
	fs := http.FileServer(http.FS(fs))
	return fs
}
