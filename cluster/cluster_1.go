package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

var addrs string
var token int

//User
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

//"localhost:9003"
var direccion string

type Info struct {
	Tipo          string
	NumNodo       int
	AddrNodo      string
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

type InfoPrediction struct {
	Tipo        string
	NumNodo     int
	AddrNodo    string
	Dni         int64
	Probability float64
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
	direccion = "localhost:9001"
	fmt.Print("Ingrese la dirección del nodo:")
	fmt.Print(direccion)

	addrs = "localhost:9002"
	fmt.Print("Host 1")
	fmt.Print(addrs)

	//2.- Generar el token
	rand.Seed(time.Now().UTC().UnixNano())

	token = rand.Intn(1000000)
	fmt.Println(token)

	//User
	/*dni = 73947420
	probability = 0
	v_age = 0.20
	v_gender = 1
	v_uci = 0
	v_oxigen = 0
	v_ventilator = 0
	v_first_dose = 0
	v_second_dose = 0
	v_vaccine = 0*/

	//sinopharm: 0.35
	//pfizer:0.65
	//astrazeneca:1

	chanIniciar = make(chan bool)
	chanMyInfo = make(chan MyInfo)

	//Establecer el valor inicial de la información del nodo
	go func() {
		chanMyInfo <- MyInfo{0, true, 1000001, ""}
	}()

	data := returnInfo()

	//3.- Iniciar el proceso
	go func() {
		fmt.Print("Presione enter para iniciar...")
		bufferIn := bufio.NewReader(os.Stdin)
		bufferIn.ReadString('\n') //pausa espera hasta q presione enter
		info := User{data.V_dni, data.V_age, data.V_gender, data.V_uci, data.V_oxigen, data.V_ventilator, data.V_first_dose, data.V_second_dose, data.V_vaccine}
		go enviar(addrs, info)

	}()

	//TODO
	//4.- Definir el servicio de acceso a la sección crítica
	ServicioSC()
}

func enviar(addr string, info User) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()
	//codificar el mensaje a enviar
	byteInfo, _ := json.Marshal(info)
	fmt.Fprintln(con, string(byteInfo))
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

func returnInfo() (info User) {
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
	var info User
	json.Unmarshal([]byte(bInfo), &info)
	fmt.Println(info)

}
