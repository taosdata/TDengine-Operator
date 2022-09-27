package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ExecResult struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func ExecSql(url string, sql, user, password string) error {
	reader := strings.NewReader(sql)
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var result ExecResult
	err = decoder.Decode(&result)
	if err != nil {
		return err
	}
	if result.Code != 0 {
		return fmt.Errorf("code: %d,desc: %s", result.Code, result.Desc)
	}
	return nil
}

type QueryResult struct {
	Code       int             `json:"code"`
	Desc       string          `json:"desc"`
	ColumnMeta [][]interface{} `json:"column_meta"`
	Data       [][]interface{} `json:"data"`
	Rows       int             `json:"rows"`
}

func GetDnodeCount(url, user, password string) (int, error) {
	result, err := Query(url, "show dnodes", user, password)
	if err != nil {
		return 0, err
	}
	return result.Rows, nil
}
func GetDnodeMap(url, user, password string) (map[string]int, error) {
	result, err := Query(url, "show dnodes", user, password)
	if err != nil {
		return nil, err
	}
	r := make(map[string]int, result.Rows)
	for _, d := range result.Data {
		r[d[1].(string)] = int(d[0].(float64))
	}
	return r, err
}

func Query(url string, sql, user, password string) (*QueryResult, error) {
	reader := strings.NewReader(sql)
	req, err := http.NewRequest(http.MethodPost, url, reader)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(user, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var result QueryResult
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("code: %d,desc: %s", result.Code, result.Desc)
	}
	return &result, nil
}
