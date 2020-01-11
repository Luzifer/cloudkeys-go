package main

import (
	"mime"
	"net/http"
	"path"
	"path/filepath"
)

func serveAssets(res http.ResponseWriter, r *http.Request) {
	var fileName = r.RequestURI[1:]
	if fileName == "" {
		fileName = "index.html"
	}

	data, err := Asset(path.Join("dist", fileName))
	if err != nil {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.RequestURI)))
	res.Write(data)
}
