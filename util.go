package main

import (
	"database/sql"
	"errors"
	"github.com/bingoohuang/go-utils"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
)

func findMerchantDataSource(tid string) (string, error) {
	db, err := sql.Open("mysql", g_dataSource)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := "SELECT DB_USERNAME, DB_PASSWORD, PROXY_IP, PROXY_PORT, DB_NAME FROM TR_F_DB WHERE MERCHANT_ID = '" + tid + "' AND STATE = '2'"
	sqlResult := go_utils.ExecuteSql(db, query, 1)
	if sqlResult.Error != nil {
		return "", sqlResult.Error
	}

	rows := len(sqlResult.Rows)
	if rows == 0 {
		return "", errors.New("no db found for tid:" + tid)
	}

	if rows > 1 {
		return "", errors.New("more than one db found")
	}

	r := sqlResult.Rows[0]
	// user:pass@tcp(127.0.0.1:3306)/db?charset=utf8
	return r[0] + ":" + r[1] + "@tcp(" + r[2] + ":" + r[3] + ")/" + r[4] + "?charset=utf8mb4,utf8&timeout=3s", nil
}

func executeTenantSql(w http.ResponseWriter, req *http.Request, querySql string) (*go_utils.ExecuteSqlResult, bool) {
	tid := strings.TrimSpace(req.FormValue("tid"))
	if tid == "" {
		http.Error(w, "tid required", 405)
		return nil, false
	}

	merchantDataSource, err := findMerchantDataSource(tid)
	if err != nil {
		http.Error(w, err.Error(), 405)
		return nil, false
	}

	db, err := sql.Open("mysql", merchantDataSource)
	if err != nil {
		http.Error(w, err.Error(), 405)
		return nil, false
	}
	defer db.Close()

	result := go_utils.ExecuteSql(db, querySql, 100)
	if result.Error != nil {
		http.Error(w, err.Error(), 405)
		return nil, false
	}

	return &result, true
}

func executeTenantSqls(w http.ResponseWriter, req *http.Request, querySqls []string) ([]go_utils.ExecuteSqlResult, bool) {
	tid := strings.TrimSpace(req.FormValue("tid"))
	if tid == "" {
		http.Error(w, "tid required", 405)
		return nil, false
	}

	merchantDataSource, err := findMerchantDataSource(tid)
	if err != nil {
		http.Error(w, err.Error(), 405)
		return nil, false
	}

	db, err := sql.Open("mysql", merchantDataSource)
	if err != nil {
		http.Error(w, err.Error(), 405)
		return nil, false
	}
	defer db.Close()

	results := make([]go_utils.ExecuteSqlResult, len(querySqls))

	for index, querySql := range querySqls {
		result := go_utils.ExecuteSql(db, querySql, 100)
		results[index] = result
		if result.Error != nil {
			http.Error(w, result.Error.Error(), 405)
			return results, false
		}

	}

	return results, true
}
