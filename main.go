package main

import (
	"flag"
	"fmt"
	"github.com/bingoohuang/go-utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	contextPath  string
	port         string
	g_dataSource *string

	devMode *bool // to disable css/js minify

	northProxy *string
	southProxy *string

	authParam go_utils.MustAuthParam
)

func init() {
	contextPathArg := flag.String("contextPath", "", "context path")
	portArg := flag.Int("port", 8979, "Port to serve.")
	g_dataSource = flag.String("dataSource", "user:pass@tcp(127.0.0.1:3306)/?charset=utf8", "dataSource string.")

	devMode = flag.Bool("devMode", false, "devMode(disable js/css minify)")

	northProxy = flag.String("northProxy", "http://127.0.0.1:8092", "northProxy")
	southProxy = flag.String("southProxy", "http://127.0.0.1:8082", "southProxy")

	go_utils.PrepareMustAuthFlag(&authParam)

	flag.Parse()

	contextPath = *contextPathArg
	if contextPath != "" && strings.Index(contextPath, "/") < 0 {
		contextPath = "/" + contextPath
	}

	port = strconv.Itoa(*portArg)
}

func main() {
	r := mux.NewRouter()

	handleFunc(r, "/", serveHome, true)
	handleFunc(r, "/favicon.ico", go_utils.ServeFavicon("res/favicon.ico", MustAsset, AssetInfo), false)
	handleFunc(r, "/searchTenants", searchTenants, false)
	handleFunc(r, "/queryCourseTypes", queryCourseTypes, false)
	handleFunc(r, "/updateCourseTypes", updateCourseTypes, false)
	http.Handle("/", r)

	fmt.Println("start to listen at ", port)
	go_utils.OpenExplorerWithContext(contextPath, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	wrap := go_utils.DumpRequest(f)
	wrap = go_utils.MustAuth(wrap, authParam)
	if requiredGzip {
		wrap = go_utils.GzipHandlerFunc(wrap)
	}

	r.HandleFunc(contextPath+path, wrap)
}
