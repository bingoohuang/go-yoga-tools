package main

import (
	"net/http"
	"strings"
	"log"
)

func saveCaptcha(w http.ResponseWriter, req *http.Request) {
	mobile := strings.TrimSpace(req.FormValue("mobile"))
	captcha := strings.TrimSpace(req.FormValue("captcha"))

	key := "captcha:" + mobile + ":/login"
	httpResult, err := httpGet(southProxy + "/setCache?keys=" + key + "&value=" + captcha + "&ttl=30s")
	httpResult1, err1 := httpGet(northProxy + "/setCache?keys=" + key + "&value=" + captcha + "&ttl=30s")

	if err != nil {
		log.Println("http get", southProxy, "fail", err.Error())
	} else {
		log.Println("http get result", string(httpResult))
	}

	if err1 != nil {
		log.Println("http get", northProxy, "fail", err.Error())
	} else {
		log.Println("http get result", string(httpResult1))
	}
}
