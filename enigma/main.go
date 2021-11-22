package main

import (
	//"bufio"
	//"encoding/json"
	"enigma/api"
	"fmt"
	//"net"
	"net/http"

	"github.com/gorilla/mux"
)
var port string = "8080"

func main()  {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.HandleFunc("/data", api.GetData)
	apiRouter.HandleFunc("/all",api.GetAll).Methods("GET")
	apiRouter.HandleFunc("/{id}",api.GetById).Methods("GET")
	apiRouter.HandleFunc("/post",api.PostData).Methods("POST")

	apiRouter.HandleFunc("/predict2",api.PostData2).Methods("POST")

	apiRouter.HandleFunc("/predict",api.AnswerPredict).Methods("POST")
	fmt.Println("Servidor ejecutandose en el puerto",port)
	//api.ServicioSC()
	http.ListenAndServe(":"+port,router)
//	ServicioSC()
}