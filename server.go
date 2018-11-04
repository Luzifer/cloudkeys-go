package main

import (
	"net/http"
)

func init() {
	initialize()
	http.Handle("/", router())
}
