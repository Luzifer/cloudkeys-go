// +build appengine

package main

import (
	"context"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func getHTTPClient(ctx context.Context) *http.Client {
	return urlfetch.Client(ctx)
}

func getContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func main() {
	initializeStorage()
	appengine.Main()
}
