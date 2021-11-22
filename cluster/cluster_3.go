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

type MyInfo struct {
	contadorMsg int
	primero     bool
	proxNum     int
	proxAddr    string
}

var chanIniciar chan bool
var chanMyInfo chan MyInfo

func main() {
	fmt.Print("Ingrese el dirección del nodo:")
	fmt.Scanf("%s\n", &direccion)

	fmt.Printf("Host %d = ", 1)
	fmt.Scanf("%s\n", &(addrs))

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
		info := Info{"ENVIOTOKEN", token, direccion, float64(data.Probability), data.V_age, float64(data.V_gender), float64(data.V_uci), float64(data.V_oxigen), float64(data.V_ventilator), float64(data.V_first_dose), float64(data.V_second_dose), float64(data.V_vaccine)}
		go enviar(addrs, info)

	}()

	//TODO
	//4.- Definir el servicio de acceso a la sección crítica
	ServicioSC()
}

func enviar(addr string, info Info) {
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

func returnInfo() (info Info) {
	ln, _ := net.Listen("tcp", direccion)
	defer ln.Close()
	con, _ := ln.Accept()

	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	json.Unmarshal([]byte(bInfo), &info)
	return
}

func manejadorConexion(con net.Conn) {
	//lógica
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	var info Info
	json.Unmarshal([]byte(bInfo), &info)
	fmt.Println(info)
}
