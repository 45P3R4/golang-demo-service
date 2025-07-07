package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
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

	var data Order

	value, isMapContainsKey := cache[id]

	//Get data
	if isMapContainsKey {
		//From cache
		data = value
	} else {
		//From db
		data, err = DbGetRowById(id.String())
		if err != nil {
			fmt.Fprintf(w, "Failed to get data ftom database: %v", err)
		}
		cache[id] = data
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w, "Failed to marshal JSON")
	}

	fmt.Fprintf(w, "%s", dataJson)
}

func main() {

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
