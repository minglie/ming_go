package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func getJSON(db *sql.DB, sqlString string) (string, error) {
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/ming?charset=utf8")
	if err != nil {
		fmt.Printf("connect mysql fail ! [%s]", err)
	} else {
		fmt.Println("connect to mysql success")
	}
	r, _ := getJSON(db, "select id,name from mi_api")
	fmt.Println(r)
}
