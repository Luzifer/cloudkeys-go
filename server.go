package main

import (
	"net/http"
)

//go:generate go-bindata assets templates

func init() {
	initialize()
	http.Handle("/", router())
}
