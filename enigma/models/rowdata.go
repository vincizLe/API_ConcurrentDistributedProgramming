package models

type RowData struct {
	Dni          int    `json:"id_persona"`
	Edad         string `json:"edad"`
	Sexo         string `json:"sexo"`
	FlagUci      string `json:"flag_uci"`
	Oxigeno      string `json:"con_oxigeno"`
	Ventilacion  string `json:"con_ventilacion"`
	FlagVacuna   string `json:"flag_vacuna"`
	FabricDosis1 string `json:"fabricante_dosis1"`
	FabricDosis2 string `json:"fabricante_dosis2"`
}
