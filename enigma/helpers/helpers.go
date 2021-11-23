package helpers

import (
	"encoding/csv"
	"enigma/models"
	"io"
	"log"
	"net/http"
	"strconv"
)

func ReadCSVFromUrl(url string) (models.Data, error) {
	var data models.Data
	csvFile, _ := http.Get(url)
	reader := csv.NewReader(csvFile.Body)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	reader.Comma = ','
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		dni, _ := strconv.Atoi(line[5])
		age, _ := strconv.ParseFloat(line[6], 64)
		gender, _ := strconv.ParseFloat(line[7], 64)
		uci, _ := strconv.ParseFloat(line[9], 64)
		oxigen, _ := strconv.ParseFloat(line[12], 64)
		ventilador, _ := strconv.ParseFloat(line[13], 64)
		dosis1, _ := strconv.ParseFloat(line[16], 64)
		dosis2, _ := strconv.ParseFloat(line[18], 64)
		vaccine, _ := strconv.ParseFloat(line[20], 64)
		data.Data = append(data.Data, models.RowData{
			Dni:          dni,
			Edad:         age,
			Sexo:         gender,
			FlagUci:      uci,
			Oxigeno:      oxigen,
			Ventilacion:  ventilador,
			FlagVacuna:   dosis1,
			FabricDosis1: dosis2,
			FabricDosis2: vaccine,
		})
	}
	return data, nil
}
