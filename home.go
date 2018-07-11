package main

import (
	"github.com/bingoohuang/go-utils"
	"net/http"
	"strings"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	indexHtml := string(MustAsset("res/index.html"))
	indexHtml = strings.Replace(indexHtml, "<LOGIN/>", loginHtml(w, r), 1)

	html := go_utils.MinifyHtml(indexHtml, *devMode)

	css := go_utils.MinifyCss(mergeCss(), *devMode)
	js := go_utils.MinifyJs(mergeScripts(), *devMode)
	html = strings.Replace(html, "/*.CSS*/", css, 1)
	html = strings.Replace(html, "/*.SCRIPT*/", js, 1)
	html = strings.Replace(html, "${contextPath}", contextPath, -1)

	w.Write([]byte(html))
}

func mergeCss() string {
	return go_utils.MergeCss(MustAsset, "index.css", "jquery.contextMenu.css")
}

func mergeScripts() string {
	return go_utils.MergeJs(MustAsset,
		"index.js", "jquery.loading.js", "courseTypes.js", "searchTenants.js")
}

func loginHtml(w http.ResponseWriter, r *http.Request) string {
	cookie := r.Context().Value("CookieValue").(*go_utils.CookieValueImpl)

	return `<img class="loginAvatar" src="` + cookie.Avatar +
		`"/><span class="loginName">` + cookie.Name + `</span>`
}
