package models

type Prediction struct {
	Dni         string `json:"id_persona"`
	Probability string `json:"probability"`
}
