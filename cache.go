package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
)

const cacheSize = 10

var DbCache *cache.Cache

func tryGetFromCache(id uuid.UUID) (Order, error) {

	var data Order

	var err error

	order, found := DbCache.Get(id.String())
	if found {
		//From cache
		data = order.(Order)
	} else {
		//From db
		data, err = DbGetRowById(id.String())
		DbCache.Set(id.String(), data, cache.DefaultExpiration)
	}

	return data, err
}

func fillCache() {
	orders, err := DbGetLastRows(cacheSize)
	if err != nil {
		fmt.Println("Failed to fill cache ", err)
	}
	for _, order := range orders {
		// cache[order.OrderUID] = order
		DbCache.Set(order.OrderUID.String(), order, 0)
	}
}
