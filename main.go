// +build !appengine

package main

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func getHTTPClient(ctx context.Context) *http.Client {
	return &http.Client{}
}

func getContext(r *http.Request) context.Context {
	return r.Context()
}

func main() {
	initializeStorage()
	log.WithError(http.ListenAndServe(cfg.Listen, nil)).Error("HTTP Server exited unexpectedly")
}
