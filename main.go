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
	contextPath       string
	port              string
	g_dataSource      string
	writeAuthRequired bool
	encryptKey        string

	corpId      string
	corpSecret  string
	agentId     string
	redirectUri string

	cookieName string
	devMode    bool // to disable css/js minify
	authBasic  bool

	northProxy string
	southProxy string
)

func init() {
	contextPathArg := flag.String("contextPath", "", "context path")
	portArg := flag.Int("port", 8979, "Port to serve.")
	dataSourceArg := flag.String("dataSource", "user:pass@tcp(127.0.0.1:3306)/?charset=utf8", "dataSource string.")
	writeAuthRequiredArg := flag.Bool("writeAuthRequired", false, "write auth required")
	encryptKeyArg := flag.String("encryptKey", "B185277FC6C4E6AB", "key to encryption or decryption")
	corpIdArg := flag.String("corpId", "", "corpId")
	corpSecretArg := flag.String("corpSecret", "", "cropId")
	agentIdArg := flag.String("agentId", "", "agentId")
	redirectUriArg := flag.String("redirectUri", "", "redirectUri")
	cookieNameArg := flag.String("cookieName", "yoga_qyapi", "cookieName")
	devModeArg := flag.Bool("devMode", false, "devMode(disable js/css minify)")
	authBasicArg := flag.Bool("authBasic", false, "authBasic based on poems")
	northProxycArg := flag.String("northProxy", "http://127.0.0.1:8092", "northProxy")
	southProxycArg := flag.String("southProxy", "http://127.0.0.1:8082", "southProxy")

	flag.Parse()

	contextPath = *contextPathArg
	if contextPath != "" && strings.Index(contextPath, "/") < 0 {
		contextPath = "/" + contextPath
	}

	port = strconv.Itoa(*portArg)
	g_dataSource = *dataSourceArg
	writeAuthRequired = *writeAuthRequiredArg
	encryptKey = *encryptKeyArg
	corpId = *corpIdArg
	corpSecret = *corpSecretArg
	agentId = *agentIdArg
	redirectUri = *redirectUriArg
	cookieName = *cookieNameArg
	devMode = *devModeArg
	authBasic = *authBasicArg

	northProxy = *northProxycArg
	southProxy = *southProxycArg
}

func main() {
	r := mux.NewRouter()

	handleFunc(r, "/", serveHome, true)
	handleFunc(r, "/favicon.ico", serveFavicon, false)
	handleFunc(r, "/searchTenants", searchTenants, false)
	handleFunc(r, "/queryCourseTypes", queryCourseTypes, false)
	handleFunc(r, "/updateCourseTypes", updateCourseTypes, false)
	if writeAuthRequired {
		handleFunc(r, "/login", serveLogin, false)
	}
	http.Handle("/", r)

	fmt.Println("start to listen at ", port)
	go_utils.OpenExplorerWithContext(contextPath, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request), requiredGzip bool) {
	wrap := go_utils.DumpRequest(f)
	if requiredGzip {
		wrap = go_utils.GzipHandlerFunc(wrap)
	}

	r.HandleFunc(contextPath+path, wrap)
}
