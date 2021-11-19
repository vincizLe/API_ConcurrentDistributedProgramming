package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type Errors struct {
	Errors []float64
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

func fit(X [][]float64, D [750]float64, Epochs int) (Errors [100]float64, o_weight [16]float64, o_bias [5]float64) {
	rows := len(X)
	Weight := [16]float64{0.1, 0.9, 0.2, 0.8, 0.3, 0.7, 0.2, 0.1, 0.3, 0.9, 0.4, 0.5, 0.3, 0.7, 0.6, 0.1}
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
			net_h1 := (Weight[0] * X[j][0]) + (Weight[1] * X[j][1]) + (Weight[2] * X[j][2]) + Bias[0]
			net_h2 := (Weight[3] * X[j][0]) + (Weight[4] * X[j][1]) + (Weight[5] * X[j][2]) + Bias[1]
			net_h3 := (Weight[6] * X[j][0]) + (Weight[7] * X[j][1]) + (Weight[8] * X[j][2]) + Bias[2]
			net_h4 := (Weight[9] * X[j][0]) + (Weight[10] * X[j][1]) + (Weight[11] * X[j][2]) + Bias[3]

			//Execute the sigmoid activation function in the first layer
			out_h1 := sigmoid(net_h1)
			out_h2 := sigmoid(net_h2)
			out_h3 := sigmoid(net_h3)
			out_h4 := sigmoid(net_h4)

			//We calculate the total net output of the output layer
			net_y := (Weight[12] * out_h1) + (Weight[13] * out_h2) + (Weight[14] * out_h3) + (Weight[15] * out_h4) + Bias[4]

			//We execute the sigmoid activation function in the output layer
			out_y := sigmoid(net_y)

			//--------------------WE CALCULATE THE TOTAL ERROR----------------------
			real_error := D[j] - out_y
			Total_Error = 0.5 * ((D[j] - out_y) * (D[j] - out_y))

			//---------------------------BACKPROPAGATION----------------------------
			//Delta rule in output layer
			delta_y := d_sigmoid(net_y) * real_error

			//We adjust the weights of the output layer
			Weight[12] = Weight[12] + (out_h1 * Learning_Factor * delta_y)
			Weight[13] = Weight[13] + (out_h2 * Learning_Factor * delta_y)
			Weight[14] = Weight[14] + (out_h3 * Learning_Factor * delta_y)
			Weight[15] = Weight[15] + (out_h4 * Learning_Factor * delta_y)

			//We adjust the bias
			Bias[4] = Bias[4] + (Learning_Factor * delta_y)

			//Delta rule in the hide layer
			delta_h1 := d_sigmoid(net_h1) * Weight[12] * delta_y
			delta_h2 := d_sigmoid(net_h2) * Weight[13] * delta_y
			delta_h3 := d_sigmoid(net_h3) * Weight[14] * delta_y
			delta_h4 := d_sigmoid(net_h4) * Weight[15] * delta_y

			//We adjust the weights of the input layer
			Weight[0] = Weight[0] + (delta_h1 * X[j][0] * Learning_Factor)
			Weight[1] = Weight[1] + (delta_h1 * X[j][1] * Learning_Factor)
			Weight[2] = Weight[2] + (delta_h1 * X[j][2] * Learning_Factor)
			Weight[3] = Weight[3] + (delta_h2 * X[j][0] * Learning_Factor)
			Weight[4] = Weight[4] + (delta_h2 * X[j][1] * Learning_Factor)
			Weight[5] = Weight[5] + (delta_h2 * X[j][2] * Learning_Factor)
			Weight[6] = Weight[6] + (delta_h3 * X[j][0] * Learning_Factor)
			Weight[7] = Weight[7] + (delta_h3 * X[j][1] * Learning_Factor)
			Weight[8] = Weight[8] + (delta_h3 * X[j][2] * Learning_Factor)
			Weight[9] = Weight[9] + (delta_h4 * X[j][0] * Learning_Factor)
			Weight[10] = Weight[10] + (delta_h4 * X[j][1] * Learning_Factor)
			Weight[11] = Weight[11] + (delta_h4 * X[j][2] * Learning_Factor)

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

func prediction(weights [16]float64, bias [5]float64, v_age float64, v_gender float64, v_uci float64) (out_y float64) {
	//-----------------------------Propagation------------------------------
	//We calculate the total net input of the hidden layer
	net_h1 := (weights[0] * v_age) + (weights[1] * v_gender) + (weights[2] * v_uci) + bias[0]
	net_h2 := (weights[3] * v_age) + (weights[4] * v_gender) + (weights[5] * v_uci) + bias[1]
	net_h3 := (weights[6] * v_age) + (weights[7] * v_gender) + (weights[8] * v_uci) + bias[2]
	net_h4 := (weights[9] * v_age) + (weights[10] * v_gender) + (weights[11] * v_uci) + bias[3]

	//Execute the sigmoid activation function in the first layer
	out_h1 := sigmoid(net_h1)
	out_h2 := sigmoid(net_h2)
	out_h3 := sigmoid(net_h3)
	out_h4 := sigmoid(net_h4)

	//We calculate the total net output of the output layer
	net_y := (weights[12] * out_h1) + (weights[13] * out_h2) + (weights[14] * out_h3) + (weights[15] * out_h4) + bias[4]

	//We execute the sigmoid activation function in the output layer
	out_y = sigmoid(net_y)

	return
}

// generate random data for line chart
func generateLineItems(errors Errors, epochs int) []opts.LineData {
	err := errors.Errors
	items := make([]opts.LineData, 0)
	for i := 0; i < epochs; i++ {
		items = append(items, opts.LineData{Value: err[i]})
	}
	return items
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	//Declare the the dimension of the data
	var slice = make([][]float64, 750)
	var D [750]float64

	//Open the file
	open_file, err := os.Open("data.csv")
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

			//Add to the array
			slice[tmp] = []float64{age / 100, gender, uci}
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
	Epochs := 100

	//Training
	e, weights, bias := fit(slice, D, Epochs)

	fmt.Println(weights, bias)

	//
	errors := Errors{Errors: e[:]}

	//Create a new line instance
	line := charts.NewLine()
	//Set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Epoch x Total Error",
			Subtitle: "Enigma",
		}))

	// Put data into instance
	var a = make([]int, Epochs)
	for i := 0; i < Epochs; i++ {
		a[i] = i
	}

	line.SetXAxis(a).
		AddSeries("Category A", generateLineItems(errors, Epochs)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)

	//Prediction
	v_age := 0.30
	v_gender := 0.00
	v_uci := 0

	out_y := prediction(weights, bias, v_age, v_gender, float64(v_uci))
	fmt.Println(out_y)
}

func main() {
	http.HandleFunc("/", httpserver)
	http.ListenAndServe(":8081", nil)
}
