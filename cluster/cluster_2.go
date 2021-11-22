package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"encoding/csv"
	"io"
	"math"
	"strconv"
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

	weights, bias := test()

	data := returnInfo()

	probability = prediction(weights, bias, data.V_age, float64(data.V_gender), float64(data.V_uci), float64(data.V_oxigen), float64(data.V_ventilator), float64(data.V_first_dose), float64(data.V_second_dose), float64(data.V_vaccine))
	//probability = prediction(weights, bias, v_age, float64(v_gender), float64(v_uci), float64(v_oxigen), float64(v_ventilator), float64(v_first_dose), float64(v_second_dose), float64(v_vaccine))
	//fmt.Println(data.V_age, data.V_gender, data.V_oxigen, data.V_uci, data.V_ventilator, data.V_first_dose, data.V_second_dose, data.V_vaccine)
	//Establecer el valor inicial de la información del nodo
	go func() {
		chanMyInfo <- MyInfo{0, true, 1000001, ""}
	}()

	//3.- Iniciar el proceso
	go func() {
		fmt.Print("Presione enter para iniciar...")
		bufferIn := bufio.NewReader(os.Stdin)
		bufferIn.ReadString('\n') //pausa espera hasta q presione enter
		info := Info{"ENVIOTOKEN", token, direccion, float64(probability), data.V_age, float64(data.V_gender), float64(data.V_uci), float64(data.V_oxigen), float64(data.V_ventilator), float64(data.V_first_dose), float64(data.V_second_dose), float64(data.V_vaccine)}
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

func sigmoid(x float64) (s float64) {
	s = 1 / (1 + math.Exp(-x))
	return
}

func d_sigmoid(x float64) (s float64) {
	s = 1 / (1 + math.Exp(-x))
	s = s * (1 - s)
	return
}

func fit(X [][]float64, D [750]float64, Epochs int) (Errors [10000]float64, o_weight [36]float64, o_bias [5]float64) {
	rows := len(X)
	Weight := [36]float64{0.1, 0.9, 0.2, 0.8, 0.3, 0.7, 0.2, 0.1, 0.3, 0.9, 0.4, 0.5, 0.3, 0.7, 0.6, 0.1, 0.85, 0.2, 0.55, 0.1, 0.65, 0.15, 0.45, 0.85, 0.3, 0.5, 0.6, 0.4, 0.1, 0.5, 0.2, 0.9, 0.35, 0.75, 0.55, 0.8}
	Learning_Factor := 0.5
	Bias := [5]float64{1, 1, 1, 1, 1}
	Epoch := 0
	var Total_Error float64

	for Epoch < Epochs {
		j := 0
		//Itirate the epochs
		for j < rows {
			//-----------------------------Propagation------------------------------
			//We calculate the total net input of the hidden layer
			net_h1 := (Weight[0] * X[j][0]) + (Weight[1] * X[j][1]) + (Weight[2] * X[j][2]) + (Weight[3] * X[j][3]) + (Weight[4] * X[j][4]) + (Weight[5] * X[j][5]) + (Weight[6] * X[j][6]) + (Weight[7] * X[j][7]) + Bias[0]
			net_h2 := (Weight[8] * X[j][0]) + (Weight[9] * X[j][1]) + (Weight[10] * X[j][2]) + (Weight[11] * X[j][3]) + (Weight[12] * X[j][4]) + (Weight[13] * X[j][5]) + (Weight[14] * X[j][6]) + (Weight[15] * X[j][7]) + Bias[1]
			net_h3 := (Weight[16] * X[j][0]) + (Weight[17] * X[j][1]) + (Weight[18] * X[j][2]) + (Weight[19] * X[j][3]) + (Weight[20] * X[j][4]) + (Weight[21] * X[j][5]) + (Weight[22] * X[j][6]) + (Weight[23] * X[j][7]) + Bias[2]
			net_h4 := (Weight[24] * X[j][0]) + (Weight[25] * X[j][1]) + (Weight[26] * X[j][2]) + (Weight[27] * X[j][3]) + (Weight[28] * X[j][4]) + (Weight[29] * X[j][5]) + (Weight[30] * X[j][6]) + (Weight[31] * X[j][7]) + Bias[3]

			//Execute the sigmoid activation function in the first layer
			out_h1 := sigmoid(net_h1)
			out_h2 := sigmoid(net_h2)
			out_h3 := sigmoid(net_h3)
			out_h4 := sigmoid(net_h4)

			//We calculate the total net output of the output layer
			net_y := (Weight[32] * out_h1) + (Weight[33] * out_h2) + (Weight[34] * out_h3) + (Weight[35] * out_h4) + Bias[4]

			//We execute the sigmoid activation function in the output layer
			out_y := sigmoid(net_y)

			//--------------------WE CALCULATE THE TOTAL ERROR----------------------
			real_error := D[j] - out_y
			Total_Error = 0.5 * ((D[j] - out_y) * (D[j] - out_y))

			//---------------------------BACKPROPAGATION----------------------------
			//Delta rule in output layer
			delta_y := d_sigmoid(net_y) * real_error

			//We adjust the weights of the output layer
			Weight[32] = Weight[32] + (out_h1 * Learning_Factor * delta_y)
			Weight[33] = Weight[33] + (out_h2 * Learning_Factor * delta_y)
			Weight[34] = Weight[34] + (out_h3 * Learning_Factor * delta_y)
			Weight[35] = Weight[35] + (out_h4 * Learning_Factor * delta_y)

			//We adjust the bias
			Bias[4] = Bias[4] + (Learning_Factor * delta_y)

			//Delta rule in the hide layer
			delta_h1 := d_sigmoid(net_h1) * Weight[32] * delta_y
			delta_h2 := d_sigmoid(net_h2) * Weight[33] * delta_y
			delta_h3 := d_sigmoid(net_h3) * Weight[34] * delta_y
			delta_h4 := d_sigmoid(net_h4) * Weight[35] * delta_y

			//We adjust the weights of the input layer
			Weight[0] = Weight[0] + (delta_h1 * X[j][0] * Learning_Factor)
			Weight[1] = Weight[1] + (delta_h1 * X[j][1] * Learning_Factor)
			Weight[2] = Weight[2] + (delta_h1 * X[j][2] * Learning_Factor)
			Weight[3] = Weight[3] + (delta_h1 * X[j][3] * Learning_Factor)
			Weight[4] = Weight[4] + (delta_h1 * X[j][4] * Learning_Factor)
			Weight[5] = Weight[5] + (delta_h1 * X[j][5] * Learning_Factor)
			Weight[6] = Weight[6] + (delta_h1 * X[j][6] * Learning_Factor)
			Weight[7] = Weight[7] + (delta_h1 * X[j][7] * Learning_Factor)
			Weight[8] = Weight[8] + (delta_h2 * X[j][0] * Learning_Factor)
			Weight[9] = Weight[9] + (delta_h2 * X[j][1] * Learning_Factor)
			Weight[10] = Weight[10] + (delta_h2 * X[j][2] * Learning_Factor)
			Weight[11] = Weight[11] + (delta_h2 * X[j][3] * Learning_Factor)
			Weight[12] = Weight[12] + (delta_h2 * X[j][4] * Learning_Factor)
			Weight[13] = Weight[13] + (delta_h2 * X[j][5] * Learning_Factor)
			Weight[14] = Weight[14] + (delta_h2 * X[j][6] * Learning_Factor)
			Weight[15] = Weight[15] + (delta_h2 * X[j][7] * Learning_Factor)
			Weight[16] = Weight[16] + (delta_h3 * X[j][0] * Learning_Factor)
			Weight[17] = Weight[17] + (delta_h3 * X[j][1] * Learning_Factor)
			Weight[18] = Weight[18] + (delta_h3 * X[j][2] * Learning_Factor)
			Weight[19] = Weight[19] + (delta_h3 * X[j][3] * Learning_Factor)
			Weight[20] = Weight[20] + (delta_h3 * X[j][4] * Learning_Factor)
			Weight[21] = Weight[21] + (delta_h3 * X[j][5] * Learning_Factor)
			Weight[22] = Weight[22] + (delta_h3 * X[j][6] * Learning_Factor)
			Weight[23] = Weight[23] + (delta_h3 * X[j][7] * Learning_Factor)
			Weight[24] = Weight[24] + (delta_h4 * X[j][0] * Learning_Factor)
			Weight[25] = Weight[25] + (delta_h4 * X[j][1] * Learning_Factor)
			Weight[26] = Weight[26] + (delta_h4 * X[j][2] * Learning_Factor)
			Weight[27] = Weight[27] + (delta_h4 * X[j][3] * Learning_Factor)
			Weight[28] = Weight[28] + (delta_h4 * X[j][4] * Learning_Factor)
			Weight[29] = Weight[29] + (delta_h4 * X[j][5] * Learning_Factor)
			Weight[30] = Weight[30] + (delta_h4 * X[j][6] * Learning_Factor)
			Weight[31] = Weight[31] + (delta_h4 * X[j][7] * Learning_Factor)

			//We adjust the bias of the hidden layer
			Bias[0] = Bias[0] + (Learning_Factor * delta_y)
			Bias[1] = Bias[1] + (Learning_Factor * delta_y)
			Bias[2] = Bias[2] + (Learning_Factor * delta_y)
			Bias[3] = Bias[3] + (Learning_Factor * delta_y)

			j += 1
		}

		Errors[Epoch] = Total_Error
		Epoch += 1
	}
	o_weight = Weight
	o_bias = Bias
	return
}

func prediction(weights [36]float64, bias [5]float64, v_age, v_gender, v_uci, v_oxigen, v_ventilator, v_first_dose, v_second_dose, v_vaccine float64) (out_y float64) {
	//-----------------------------Propagation------------------------------
	//We calculate the total net input of the hidden layer
	net_h1 := (weights[0] * v_age) + (weights[1] * v_gender) + (weights[2] * v_uci) + (weights[3] * v_oxigen) + (weights[4] * v_ventilator) + (weights[5] * v_first_dose) + (weights[6] * v_second_dose) + (weights[7] * v_vaccine) + bias[0]
	net_h2 := (weights[8] * v_age) + (weights[9] * v_gender) + (weights[10] * v_uci) + (weights[11] * v_oxigen) + (weights[12] * v_ventilator) + (weights[13] * v_first_dose) + (weights[14] * v_second_dose) + (weights[15] * v_vaccine) + bias[1]
	net_h3 := (weights[16] * v_age) + (weights[17] * v_gender) + (weights[18] * v_uci) + (weights[19] * v_oxigen) + (weights[20] * v_ventilator) + (weights[21] * v_first_dose) + (weights[22] * v_second_dose) + (weights[23] * v_vaccine) + bias[2]
	net_h4 := (weights[24] * v_age) + (weights[25] * v_gender) + (weights[26] * v_uci) + (weights[27] * v_oxigen) + (weights[28] * v_ventilator) + (weights[29] * v_first_dose) + (weights[30] * v_second_dose) + (weights[31] * v_vaccine) + bias[3]

	//Execute the sigmoid activation function in the first layer
	out_h1 := sigmoid(net_h1)
	out_h2 := sigmoid(net_h2)
	out_h3 := sigmoid(net_h3)
	out_h4 := sigmoid(net_h4)

	//We calculate the total net output of the output layer
	net_y := (weights[32] * out_h1) + (weights[33] * out_h2) + (weights[34] * out_h3) + (weights[35] * out_h4) + bias[4]

	//We execute the sigmoid activation function in the output layer
	out_y = sigmoid(net_y)

	return
}

func test() (weights [36]float64, bias [5]float64) {
	//Declare the the dimension of the data
	var slice = make([][]float64, 750)
	var D [750]float64

	//Open the file
	open_file, err := os.Open("data_1.csv")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened CSV file")
	}

	//Read the file
	read_file := csv.NewReader(open_file)

	tmp := 0
	for {
		data, err := read_file.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		} else {
			//Age
			age, err := strconv.ParseFloat(data[0], 64)
			if err != nil {
				fmt.Println(err)
			}
			//Gender
			var gender float64
			if data[1] == "F" {
				gender = 1
			} else {
				gender = 0
			}
			//UCI
			uci, err := strconv.ParseFloat(data[2], 64)
			if err != nil {
				fmt.Println(err)
			}
			//With oxigen
			oxigen, err := strconv.ParseFloat(data[3], 64)
			if err != nil {
				fmt.Println(err)
			}
			//With ventilation
			ventilation, err := strconv.ParseFloat(data[4], 64)
			if err != nil {
				fmt.Println(err)
			}
			//1st dose
			var first_dose float64
			if data[6] == "0" {
				first_dose = 0
			} else {
				first_dose = 1
			}
			//2st dose
			var second_dose float64
			if data[7] == "0" {
				second_dose = 0
			} else {
				second_dose = 1
			}
			//Vaccine
			vaccine, err := strconv.ParseFloat(data[8], 64)
			if err != nil {
				fmt.Println(err)
			}

			//Add to the array
			slice[tmp] = []float64{age / 100, gender, uci, oxigen, ventilation, first_dose, second_dose, vaccine}
			//fmt.Println(slice)
			//Output
			door := data[5] == "alta"
			if door == true {
				D[tmp] = 1

			} else {
				D[tmp] = 0
			}
			tmp = tmp + 1
		}
	}
	//Epochs
	Epochs := 10000

	//Training
	_, Weights, Bias := fit(slice, D, Epochs)

	weights = Weights
	bias = Bias
	return
}
