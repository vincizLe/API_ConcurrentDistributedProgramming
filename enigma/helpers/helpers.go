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
		dni,_ := strconv.Atoi(line[5])
		data.Data = append(data.Data, models.RowData{
			Dni:			dni,
			Edad:			line[6],
			Sexo:			line[7],
			FlagUci:		line[9],
			Oxigeno:		line[12],
			Ventilacion:	line[13],
			FlagVacuna:		line[16],
			FabricDosis1:	line[18],
			FabricDosis2:	line[20],
		})
	}
	return data, nil
}
