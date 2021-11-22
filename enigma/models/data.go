package models

type Data struct {
	Success bool      `json:"success"`
	Data    []RowData `json:"data"`
	Errors  []string  `json:"errors"`
}

type Answer struct {
	Success bool         `json:"success"`
	Answer  []Prediction `json:"answer"`
	Errors  []string     `json:"errors"`
}
