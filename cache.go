package main

import (
	"fmt"

	"github.com/google/uuid"
)

const cacheSize = 10

var cache = map[uuid.UUID]Order{}

func fillCache() {
	orders, err := DbGetLastRows(cacheSize)
	if err != nil {
		fmt.Println("Failed to fill cache ", err)
	}
	for _, order := range orders {
		cache[order.OrderUID] = order
	}
}
