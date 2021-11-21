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
var age float64
var gender float64

//"localhost:9003"
var direccion string

type Info struct {
	Tipo     string
	NumNodo  int
	AddrNodo string
	Age      float64
	Gender   float64
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

	//User
	age = 32
	gender = 35

	chanIniciar = make(chan bool)
	chanMyInfo = make(chan MyInfo)

	//Establecer el valor inicial de la información del nodo
	go func() {
		chanMyInfo <- MyInfo{0, true, 1000001, ""}
	}()

	//3.- Iniciar el proceso
	go func() {
		fmt.Print("Presione enter para iniciar...")
		bufferIn := bufio.NewReader(os.Stdin)
		bufferIn.ReadString('\n') //pausa espera hasta q presione enter
		info := Info{"ENVIOTOKEN", token, direccion, age, gender}
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

func manejadorConexion(con net.Conn) {
	//lógica
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	var info Info
	json.Unmarshal([]byte(bInfo), &info)
	fmt.Println(info)

	//Evaluar según el tipo de mensaje
	switch info.Tipo {
	case "ENVIOTOKEN":
		//recuperar del canal la info del nodo
		myInfo := <-chanMyInfo
		if info.NumNodo < token {
			myInfo.primero = false
		} else if info.NumNodo < myInfo.proxNum {
			myInfo.proxAddr = info.AddrNodo
			myInfo.proxNum = info.NumNodo
		}
		//actualiza el numero de nodos notificados
		myInfo.contadorMsg++
		//retorno por canal con la info actual
		go func() {
			chanMyInfo <- myInfo
		}()
		//evaluar el fin del proceso
		if myInfo.contadorMsg == len(addrs) {
			//evaluar
			if myInfo.primero {
				procesarSC()
			} else {
				chanIniciar <- true //sincronización, pausa
			}
		}
	case "INICIO":
		<-chanIniciar //espera hasta que llegue true
		procesarSC()
	}
}

func procesarSC() {
	fmt.Println("Inicia el proceso")
	//evalua el proximo a procesar y si es el último
	myInfo := <-chanMyInfo
	if myInfo.proxAddr == "" {
		fmt.Println("Soy el último nodo en procesar!")

	} else {
		fmt.Println("No soy el último, procesando SC!!")
		fmt.Println(myInfo.proxAddr)
		//envia notificación al próximo nodo a procesar su SC
		info := Info{Tipo: "INICIO"}
		enviar(myInfo.proxAddr, info) //próximo nodo
	}
}
