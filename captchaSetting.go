package main

import (
	"net/http"
	"strings"
	"github.com/bingoohuang/go-utils"
)

func saveCaptcha(w http.ResponseWriter, req *http.Request) {
	mobile := strings.TrimSpace(req.FormValue("mobile"))
	activeHomeArea := strings.TrimSpace(req.FormValue("activeHomeArea"))
	activeClassifier := strings.TrimSpace(req.FormValue("activeClassifier"))
	proxy := getProxy(activeHomeArea)

	key := ""
	if activeClassifier == "et" {
		key = "captcha:" + mobile + ":/login/sms"
	} else {
		key = "captcha:" + mobile + ":/login"
	}
	setCacheUrl := proxy + "/setCache?key=" + key + "&value=1234&ttl=60s"

	httpResult, err := go_utils.HttpGet(setCacheUrl)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if err != nil {
		http.Error(w, err.Error(), 405)
	} else {
		w.Write(httpResult)
	}

}
