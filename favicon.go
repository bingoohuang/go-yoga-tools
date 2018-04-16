package main

import (
	"net/http"
	"bytes"
	"github.com/bingoohuang/go-utils"
	"io"
)

func serveFavicon(w http.ResponseWriter, r *http.Request) {
	HandleStaticResource("res/favicon.ico", w)
}

func HandleStaticResource(path string, w http.ResponseWriter) {
	data := MustAsset(path)
	fi, _ := AssetInfo(path)
	buffer := bytes.NewReader(data)
	w.Header().Set("Content-Type", go_utils.DetectContentType(fi.Name()))
	w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)
	io.Copy(w, buffer)
}


