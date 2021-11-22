package api

import (
	"bufio"
	"encoding/json"
	"enigma/helpers"
	"enigma/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)


type DataPredict struct{
	Dni			string `json:"id_persona"`
	Percentage	string `json:"predict_percent"`
}

var v_dni int
var v_age float64
var v_gender float64
var v_uci float64
var v_oxigen float64
var v_ventilator float64
var v_first_dose float64
var v_second_dose float64
var v_vaccine float64

type User struct{
	V_dni		  int
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
	for _, obj := range data.Data{
		dni, _ := strconv.Atoi(params["id"])
		if obj.Dni == dni {
			json.NewEncoder(res).Encode(obj)
			return
		}
	}
	json.NewEncoder(res).Encode(&data.Errors)
}

/*func ServicioSC() {
	ln, _ := net.Listen("tcp","localhost:8080")
	defer ln.Close()
	fmt.Println("funciona1")
	for {
		con, _ := ln.Accept()
		fmt.Println("funciona2")
		go manejadorConexion(con)
		fmt.Println("funciona3")
	}

}*/

/*func manejadorConexion(con net.Conn) {
	//lógica
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	var user User
	json.Unmarshal([]byte(bInfo), &user)
	fmt.Println(user)
	
}*/

func PostData(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	var data2 models.RowData
	_ = json.NewDecoder(req.Body).Decode(&data2)
	data.Data = append(data.Data, data2)
	json.NewEncoder(res).Encode(data.Data)
}
//------------------------------------------------------------------
func PostData2(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	var data2 models.RowData

	_ = json.NewDecoder(req.Body).Decode(&data2)
	data.Data = append(data.Data, data2)
	v_dni 		= data2.Dni
	v_age 		= 1//data2.Edad
	v_gender	=	1//data2.Sexo 
	v_uci		=   1//data2.FlagUci
	v_oxigen	=	1//data2.Oxigeno
	v_ventilator=	1//data2.Ventilacion
	v_first_dose =	1//data2.FabricDosis1
	v_second_dose=	1//data2.FabricDosis2
	v_vaccine	=	1//data2.FlagVacuna

	json.NewEncoder(res).Encode(data.Data)

	user := User{ v_dni, v_age, v_gender, v_uci, v_oxigen, v_ventilator, v_first_dose, v_second_dose, v_vaccine}
	//go enviar("localhost:8080", user)
	//ServicioSC()
	fmt.Println(user,"gaaa")
}

/*func enviar(addr string, user User) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()
	//codificar el mensaje a enviar
	byteUser, _ := json.Marshal(user)
	fmt.Fprintln(con, string(byteUser))
	fmt.Println("waaa")
	fmt.Println(user)
}*/
	
//--------------------------------
func AnswerPredict(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == "POST"{
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
			http.Error(res, "Error en body", http.StatusInternalServerError)
		}
		var data2 models.RowData
		//_ = json.NewDecoder(req.body).Decode(&data2)
		/*data_ejemplo := DataPredict{
			Dni: 		"11111",
			Percentage: "0.545646",
		}*/
		json.Unmarshal(body, &data2)
		SendPredict(data2)
		
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		io.WriteString(res, `{
			"msg":"Registro Data Predicicón correcta"
		}`)
	} else {
		http.Error(res, "Método no válido", http.StatusMethodNotAllowed)
	}
	
}

func SendPredict(msg models.RowData){
	con,_ := net.Dial("tcp", "localhost:8080/api/predict2")

	bytmsg,_ := json.Marshal(msg)
	/*
	if err != nil {
		fmt.Fprint(con,"Falla en la identificacion de la direccion!",err.Error())
	}*/
	fmt.Fprintf(con, "La prediccion es:%s\n", string(bytmsg))

	//msg1, _ := bufio.NewReader(conn).ReadString('\n')

	bufferIn := bufio.NewReader(con)  //buffer de entrada
	msg1, _ := bufferIn.ReadString('\n')

	fmt.Println("Mensaje del servidor: "+msg1)
	con.Close()
}

