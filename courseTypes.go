package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"github.com/bingoohuang/go-utils"
	"net/url"
	"io/ioutil"
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
	_, ok := executeTenantSqls(w, req, sqls)
	if ok {
		// clear redis cache
		proxy := ""
		if activeHomeArea == "south-center" {
			proxy = northProxy
		} else if activeHomeArea == "north-center" {
			proxy = southProxy
		}

		httpGet(proxy + "/clearCache?key=" + url.QueryEscape("westcache:yoga:"+tcode+":CourseTypeDaoImpl.queryCourseTypes"))
	}
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
