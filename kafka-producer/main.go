package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// var testData = []byte(`{
// 		"order_uid": "b563feb7b2b84b6test",
// 		"track_number": "WBILMTESTTRACK",
// 		"entry": "WBIL",
// 		"delivery": {
// 			"delivery_uid": "b563feb7b2b84b6test",
// 			"name": "Test Testov",
// 			"phone": "+9720000000",
// 			"zip": 2639809,
// 			"city": "Kiryat Mozkin",
// 			"address": "Ploshad Mira 15",
// 			"region": "Kraiot",
// 			"email": "test@gmail.com"
// 		},
// 		"payment": {
// 			"transaction": "b563feb7b2b84b6test",
// 			"request_id": "",
// 			"currency": "USD",
// 			"provider": "wbpay",
// 			"amount": 1817,
// 			"payment_dt": 1637907727,
// 			"bank": "alpha",
// 			"delivery_cost": 1500,
// 			"goods_total": 317,
// 			"custom_fee": 0
// 		},
// 		"items": [
// 			{
// 				"chrt_id": 9934930,
// 				"track_number": "WBILMTESTTRACK",
// 				"price": 453,
// 				"rid": "ab4219087a764ae0btest",
// 				"name": "Mascaras",
// 				"sale": 30,
// 				"size": 0,
// 				"total_price": 317,
// 				"nm_id": 2389212,
// 				"brand": "Vivienne Sabo",
// 				"status": 202
// 			}
// 		],
// 		"locale": "en",
// 		"internal_signature": "",
// 		"customer_id": "test",
// 		"delivery_service": "meest",
// 		"shardkey": 9,
// 		"sm_id": 99,
// 		"date_created": "2021-11-26T06:22:19Z",
// 		"oof_shard": 1
// 	}`)

var testData = []byte(`{
		"order_uid": "ba14feb9b2b84f9test",
		"track_number": "WBILMTESTTRACK",
		"entry": "WBIL",
		"delivery": {
			"delivery_uid": "ba14feb9b2b84f9test",
			"name": "Tes Tes",
			"phone": "+9080000000",
			"zip": 3731409,
			"city": "Bobrov",
			"address": "Ploshad Bobrov 29",
			"region": "Bober",
			"email": "test@gmail.com"
		},
		"payment": {
			"transaction": "ba14feb9b2b84f9test",
			"request_id": "",
			"currency": "RUB",
			"provider": "wbpay",
			"amount": 1817,
			"payment_dt": 1637907727,
			"bank": "alpha",
			"delivery_cost": 1500,
			"goods_total": 317,
			"custom_fee": 0
		},
		"items": [
			{
				"chrt_id": 9934930,
				"track_number": "WBILMTESTTRACK",
				"price": 453,
				"rid": "te9019087a824ae0btest",
				"name": "TATATOTO",
				"sale": 30,
				"size": 0,
				"total_price": 317,
				"nm_id": 2389212,
				"brand": "KEKEGOGO",
				"status": 202
			},
			{
				"chrt_id": 4214936,
				"track_number": "WBILMTESTTRACK",
				"price": 33,
				"rid": "qr4219087a574ae0btest",
				"name": "Neko Ark",
				"sale": 10,
				"size": 0,
				"total_price": 30,
				"nm_id": 2389212,
				"brand": "meme",
				"status": 202
			},
			{
				"chrt_id": 2934432,
				"track_number": "WBILMTESTTRACK",
				"price": 1500,
				"rid": "vf4219087a254ae0btest",
				"name": "Thing",
				"sale": 0,
				"size": 0,
				"total_price": 1500,
				"nm_id": 2389212,
				"brand": "Thinger",
				"status": 202
			}
		],
		"locale": "en",
		"internal_signature": "",
		"customer_id": "test",
		"delivery_service": "meest",
		"shardkey": 9,
		"sm_id": 99,
		"date_created": "2021-11-26T06:22:19Z",
		"oof_shard": 1
	}`)

func connect_kafka(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	return conn, err
}

func produceTestData() {
	var conn, err = connect_kafka("orders", 0)
	if err != nil {
		fmt.Println("Failed to connect Kafka")
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	_, err = conn.WriteMessages(
		kafka.Message{Value: testData},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	fmt.Println("message created")

}

func main() {

	produceTestData()

}
