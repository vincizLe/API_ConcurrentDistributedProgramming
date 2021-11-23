package models

type RowData struct {
	Dni          int     `json:"id_persona"`
	Edad         float64 `json:"edad"`
	Sexo         float64 `json:"sexo"`
	FlagUci      float64 `json:"flag_uci"`
	Oxigeno      float64 `json:"con_oxigeno"`
	Ventilacion  float64 `json:"con_ventilacion"`
	FlagVacuna   float64 `json:"flag_vacuna"`
	FabricDosis1 float64 `json:"fabricante_dosis1"`
	FabricDosis2 float64 `json:"fabricante_dosis2"`
}
