package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// var cache map[string]Order

func getOrderById(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	id := req.PathValue("id")

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(w, "Order not found")
		}
	}()

	data, err := DbGetRowById(id)
	if err != nil {
		fmt.Fprintf(w, "Failed to get data ftom database: %v", err)
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(w, "Failed to marshal JSON")
	}
	w.Header().Set("Content-Type", "application/json")

	fmt.Fprintf(w, "%s", dataJson)
}

func main() {

	go kafkaListen()

	http.HandleFunc("/order/{id}", getOrderById)

	http.ListenAndServe(":8081", nil)

}
