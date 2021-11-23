package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

var addrs string
var token int

//User
var probability float64
var v_age float64
var v_gender float64
var v_uci float64
var v_oxigen float64
var v_ventilator float64
var v_first_dose float64
var v_second_dose float64
var v_vaccine float64

//"localhost:9003"
var direccion string

type Info struct {
	Tipo          string
	NumNodo       int
	AddrNodo      string
	Dni           int64
	Probability   float64
	V_age         float64
	V_gender      float64
	V_uci         float64
	V_oxigen      float64
	V_ventilator  float64
	V_first_dose  float64
	V_second_dose float64
	V_vaccine     float64
}

type Prediction struct {
	Dni         string `json:"id_persona"`
	Probability string `json:"probability"`
}

type MyInfo struct {
	contadorMsg int
	primero     bool
	proxNum     int
	proxAddr    string
}

var chanIniciar chan bool
var chanMyInfo chan MyInfo

func main() {

	addrs = " localhost:8080/api/predict"
	fmt.Print("Ingrese la dirección del nodo: ")
	fmt.Scanf("%s\n", &direccion)
	fmt.Printf("Host %d = ", 1)
	fmt.Printf(addrs)

	//2.- Generar el token
	rand.Seed(time.Now().UTC().UnixNano())

	token = rand.Intn(1000000)
	fmt.Println(token)

	chanIniciar = make(chan bool)
	chanMyInfo = make(chan MyInfo)

	data := returnInfo()

	//Establecer el valor inicial de la información del nodo
	go func() {
		chanMyInfo <- MyInfo{0, true, 1000001, ""}
	}()

	//3.- Iniciar el proceso
	go func() {
		fmt.Print("Presione enter para iniciar...")
		bufferIn := bufio.NewReader(os.Stdin)
		bufferIn.ReadString('\n') //pausa espera hasta q presione enter
		info := Prediction{data.Dni, data.Probability}
		go enviar(addrs, info)

	}()

	//TODO
	//4.- Definir el servicio de acceso a la sección crítica
	ServicioSC()
}

func enviar(addr string, info Prediction) {
	prediccion := map[string]string{"id_persona": info.Dni, "probability": info.Probability}

	json_data, err := json.Marshal(prediccion)

	resp, err := http.Post("http://localhost:8080/api/predict", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["json"])

}

/////////////////////
func ServicioSC() {
	ln, _ := net.Listen("tcp", direccion)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejadorConexion(con)
	}
}

func returnInfo() (info Prediction) {
	ln, _ := net.Listen("tcp", direccion)
	defer ln.Close()
	con, _ := ln.Accept()

	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	json.Unmarshal([]byte(bInfo), &info)
	fmt.Println(info)

	return
}

func manejadorConexion(con net.Conn) {
	//lógica
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	var info Prediction
	json.Unmarshal([]byte(bInfo), &info)
	fmt.Println(info)
}
