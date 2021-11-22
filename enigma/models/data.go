package models

type Data struct {
	Success bool `json:"success"`
	Data []RowData `json:"data"`
	Errors []string `json:"errors"`
}