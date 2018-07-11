package main

import (
	"database/sql"
	"encoding/json"
	"github.com/bingoohuang/go-utils"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

type SearchResult struct {
	MerchantName string
	MerchantId   string
	MerchantCode string
	HomeArea     string
	Classifier   string
}

func searchTenants(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	searchKey := strings.TrimSpace(req.FormValue("tenantSearchKey"))
	if searchKey == "" {
		http.Error(w, "tenantSearchKey required", 405)
		return
	}

	searchSql := "SELECT MERCHANT_NAME, MERCHANT_ID, MERCHANT_CODE, HOME_AREA, CLASSIFIER " +
		"FROM TR_F_MERCHANT WHERE MERCHANT_ID = '" + searchKey +
		"' OR MERCHANT_CODE = '" + searchKey + "' OR MERCHANT_NAME LIKE '%" + searchKey + "%'"

	db, err := sql.Open("mysql", *g_dataSource)
	if err != nil {
		http.Error(w, err.Error(), 405)
		return
	}
	defer db.Close()

	result := go_utils.ExecuteSql(db, searchSql, 100)
	if result.Error != nil {
		http.Error(w, err.Error(), 405)
		return
	}

	searchResult := make([]SearchResult, len(result.Rows))
	for i, v := range result.Rows {
		searchResult[i] = SearchResult{
			MerchantName: v[0], MerchantId: v[1], MerchantCode: v[2], HomeArea: v[3], Classifier: v[4]}
	}

	json.NewEncoder(w).Encode(searchResult)
}
