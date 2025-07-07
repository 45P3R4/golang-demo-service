package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

func getOrderById(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //Temp solution, unknown context
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow", "GET")

	//Error handling
	defer func() {
		if r := recover(); r != nil {
			file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatal("Failed to open log file:", err)
			}
			log.SetOutput(file)
			log.Println("ERROR: ", r)
		}
	}()

	//id string to UUID
	id, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		fmt.Println("Failed to parse UUID: ", err)
	}

	data, err := tryGetFromCache(id)
	if err != nil {
		panic("[getOrderById]: Failed to get data: " + err.Error())
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w, "Failed to marshal JSON")
	}

	fmt.Fprintf(w, "%s", dataJson)
}

func main() {

	DbCache = cache.New(5*time.Minute, 10*time.Minute)
	//Error handling
	defer func() {
		if r := recover(); r != nil {
			file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatal("Failed to open log file:", err)
			}
			log.SetOutput(file)
			log.Println("ERROR: ", r)
		}
	}()

	fillCache()

	go kafkaListen()

	http.HandleFunc("/order/{id}", getOrderById)

	http.ListenAndServe(":8081", nil)

}
