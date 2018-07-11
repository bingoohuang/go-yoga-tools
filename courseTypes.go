package main

import (
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type CourseType struct {
	ShowOrder        string
	CourseTypeId     string
	CourseTypeName   string
	SubscribeType    string // 1：排的课；2：约的课
	MinSubscriptions string
}

func updateCourseTypes(w http.ResponseWriter, req *http.Request) {
	courseTypesSqls := strings.TrimSpace(req.FormValue("courseTypesSqls"))
	tcode := strings.TrimSpace(req.FormValue("tcode"))
	activeHomeArea := strings.TrimSpace(req.FormValue("activeHomeArea"))
	sqls := go_utils.SplitSqls(courseTypesSqls, ';')

	executeTenantSqls(w, req, sqls)

	// clear redis cache
	proxy := getProxy(activeHomeArea)

	if proxy == "" {
		log.Println("unknown proxy for home area", activeHomeArea)
		return
	}

	keys := strings.Replace("westcache:yoga:{tcode}:CourseTypeDaoImpl.queryCourseTypes,"+
		"westcache:yoga:{tcode}:MerchantConfigBoService.getMerchantConfig", "{tcode}", tcode, -1)
	log.Println("clear cache for keys", keys)
	httpResult, err := httpGet(proxy + "/clearCache?keys=" + url.QueryEscape(keys))
	if err != nil {
		log.Println("http get", proxy, "fail", err.Error())
	} else {
		log.Println("http get result", string(httpResult))
	}
}

func getProxy(activeHomeArea string) string {
	proxy := ""
	if strings.Contains(activeHomeArea, "south") {
		proxy = *southProxy
	} else if strings.Contains(activeHomeArea, "north") {
		proxy = *northProxy
	} else if strings.Contains(activeHomeArea, "huanan") {
		proxy = *huananProxy
	}
	return proxy
}

func httpGet(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func queryCourseTypes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	querySql := "select SHOW_ORDER, COURSE_TYPE_ID, COURSE_TYPE_NAME, SUBSCRIBE_TYPE, MIN_SUBSCRIPTIONS from tt_d_course_type where STATE = 1 order by SHOW_ORDER"

	result, ok := executeTenantSql(w, req, querySql)
	if !ok {
		return
	}

	searchResult := make([]CourseType, len(result.Rows))
	for i, v := range result.Rows {
		searchResult[i] = CourseType{
			ShowOrder: v[0], CourseTypeId: v[1], CourseTypeName: v[2], SubscribeType: v[3], MinSubscriptions: v[4]}
	}

	json.NewEncoder(w).Encode(searchResult)
}
