// +build !appengine

package main // import "github.com/Luzifer/cloudkeys-go"

//go:generate go-bindata assets

import (
	"context"
	"net/http"
)

func getHTTPClient(ctx context.Context) *http.Client {
	return &http.Client{}
}

func getContext(r *http.Request) context.Context {
	return r.Context()
}

func main() {
	initializeStorage()
	http.ListenAndServe(cfg.Listen, nil)
}
