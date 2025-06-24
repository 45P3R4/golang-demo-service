package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "orders\n")
}

func main() {

	http.HandleFunc("/order", hello)
	// http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8081", nil)
}
