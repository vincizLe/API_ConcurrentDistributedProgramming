package api

import (
	"encoding/json"
	"enigma/helpers"
	"enigma/models"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var v_dni int
var probability float64
var v_age float64
var v_gender float64
var v_uci float64
var v_oxigen float64
var v_ventilator float64
var v_first_dose float64
var v_second_dose float64
var v_vaccine float64

type User struct {
	V_dni         int
	V_age         float64
	V_gender      float64
	V_uci         float64
	V_oxigen      float64
	V_ventilator  float64
	V_first_dose  float64
	V_second_dose float64
	V_vaccine     float64
}

var data models.Data
var answer models.Answer

func GetData(res http.ResponseWriter, req *http.Request) {
	//var data models.Data
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	data, _ = helpers.ReadCSVFromUrl("https://raw.githubusercontent.com/rcucho/pueba_api/main/dataset.csv")
	res.Header().Set("Content-Type", "application/json")
	jsonBytes, _ := json.MarshalIndent(data.Data[1:], "", " ")
	io.WriteString(res, string(jsonBytes))
}

func GetAll(res http.ResponseWriter, req *http.Request) {
	data, _ = helpers.ReadCSVFromUrl("https://raw.githubusercontent.com/rcucho/pueba_api/main/dataset.csv")
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(data.Data[1:])
}

func GetById(res http.ResponseWriter, req *http.Request) {
	//var data models.Data
	data, _ = helpers.ReadCSVFromUrl("https://raw.githubusercontent.com/rcucho/pueba_api/main/dataset.csv")
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	//ServicioSC()
	for _, obj := range data.Data {
		dni, _ := strconv.Atoi(params["id"])
		if obj.Dni == dni {
			json.NewEncoder(res).Encode(obj)
			return
		}
	}
	json.NewEncoder(res).Encode(&data.Errors)
}

func PostData(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var data2 models.RowData
	_ = json.NewDecoder(req.Body).Decode(&data2)
	data.Data = append(data.Data, data2)
	json.NewEncoder(res).Encode(data.Data)
}

//------------------------------------------------------------------
func PostQuery(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var data2 models.RowData

	_ = json.NewDecoder(req.Body).Decode(&data2)
	data.Data = append(data.Data, data2)
	v_dni = data2.Dni
	v_age = 1         //data2.Edad
	v_gender = 1      //data2.Sexo
	v_uci = 1         //data2.FlagUci
	v_oxigen = 1      //data2.Oxigeno
	v_ventilator = 1  //data2.Ventilacion
	v_first_dose = 1  //data2.FabricDosis1
	v_second_dose = 1 //data2.FabricDosis2
	v_vaccine = 1     //data2.FlagVacuna

	json.NewEncoder(res).Encode(data.Data)

	user := User{v_dni, v_age, v_gender, v_uci, v_oxigen, v_ventilator, v_first_dose, v_second_dose, v_vaccine}
	go enviar("localhost:9001", user)
	//ServicioSC()
	fmt.Println(user, "gaaa")
}

func enviar(addr string, user User) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()
	//codificar el mensaje a enviar
	byteUser, _ := json.Marshal(user)
	fmt.Fprintln(con, string(byteUser))

}

func GetPredictByDNI(res http.ResponseWriter, req *http.Request) {
	var prediction models.Answer
	params := mux.Vars(req)
	//ServicioSC()
	for _, obj := range answer.Answer {
		dni, _ := params["id_persona"]
		if obj.Dni == dni {
			json.NewEncoder(res).Encode(obj)
			return
		}
	}
	json.NewEncoder(res).Encode(&prediction.Errors)
}

//--------------------------------
func PostPredict(res http.ResponseWriter, req *http.Request) {
	enableCors(&res)
	res.Header().Set("Content-Type", "application/json")

	var prediction models.Prediction

	_ = json.NewDecoder(req.Body).Decode(&prediction)
	answer.Answer = append(answer.Answer, prediction)

	json.NewEncoder(res).Encode(answer.Answer)

}
func enableCors(res *http.ResponseWriter) {
	(*res).Header().Set("Access-Control-Allow-Origin", "*")
	(*res).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*res).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
