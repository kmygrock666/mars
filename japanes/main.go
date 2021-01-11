package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var myEnv MyEnv
var mysql Mysql

func init() {
	envFile := flag.String("e", "", "env file path")
	flag.Parse()

	err := myEnv.Init(*envFile)
	if err != nil {
		panic(err)
	}

	mysql.SetConfig(myEnv)
	// myRedis.SetConfig(myEnv)
	// InitDB()
	// SQLITE.create()
	log.Printf("ENV: %+v\n", myEnv)
}

func main() {
	//db connect
	timeCost, mysqlReadErr := mysql.CreateReadConnection()
	if mysqlReadErr != nil {
		panic(mysqlReadErr)
	}
	defer mysql.CloseReadConnection()
	log.Printf("Create mysql read connection cost: %+v\n", timeCost)

	timeCost, mysqlWriteErr := mysql.CreateWriteConnection()
	if mysqlWriteErr != nil {
		panic(mysqlWriteErr)
	}
	defer mysql.CloseWriteConnection()
	log.Printf("Create mysql write connection cost: %+v\n", timeCost)

	var apiV1 APIv1
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/ping", apiV1.Ping).Methods("GET")
	api.HandleFunc("/start", apiV1.LearnAccent).Methods("GET")
	api.HandleFunc("/send", apiV1.CheckAnswer).Methods("POST")

	listenAddr := fmt.Sprintf("%s:%d", myEnv.HTTPListenHost, myEnv.HTTPListenPort)
	srv := &http.Server{
		Handler:      (r),
		Addr:         listenAddr,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}
	log.Printf("Starting listen port on 'http://%s', USE Ctrl + C to stop.", listenAddr)
	httpListenErr := srv.ListenAndServe()
	panic(httpListenErr)
}
