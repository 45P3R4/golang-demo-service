package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func kafkaListen() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "orders",
		GroupID:   "consumer-orders",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}
		if err := r.CommitMessages(context.Background(), m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}

		DbInsert(m)

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func DbInsert(m kafka.Message) {
	var dataOrder Order

	dataOrder.ItemsRID = make([]string, 0)
	dataOrder.Items = make([]Item, 1)

	err := json.Unmarshal(m.Value, &dataOrder)
	if err != nil {
		log.Fatal("FAILED TO UNMARSHAL JSON: ", err)
	}

	DbInsertItems(dataOrder.Items)
	DbInsertDeliveries(dataOrder.Delivery)
	DbInsertPayments(dataOrder.Payment)

	dataOrder.DeliveryUID = dataOrder.Delivery.DeliveryUID
	dataOrder.PaymentTransaction = dataOrder.Payment.Transaction

	for _, item := range dataOrder.Items {
		dataOrder.ItemsRID = append(dataOrder.ItemsRID, item.RID)
	}

	DbInsertOrders(dataOrder)
}
