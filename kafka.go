package main

import (
	"context"
	"log"
	"os"

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

	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}
		if err := r.CommitMessages(context.Background(), m); err != nil {
			panic("[kafkaListen]: Failed to commit message: " + err.Error())
		}

		DbInsert(m)
	}

	if err := r.Close(); err != nil {
		panic("[kafkaListen]: failed to close reader: " + err.Error())
	}
}
