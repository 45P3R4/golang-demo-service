package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

const postgresUrl = "postgresql://localhost/WBORDERS?user=wbdev&password=admin"

func DbConnect(connectionUrl string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

func DbInsertOrders(order Order) {
	conn := DbConnect(postgresUrl)
	defer conn.Close(context.Background())

	_, err := conn.Exec(context.Background(),
		`INSERT INTO orders (
			order_uid, 
			track_number, 
			entry, 
			delivery_uid,
			payment_transaction, 
			items_rid, 
			locale,
			internal_signature, 
			customer_id, 
			delivery_service,
			shardkey, sm_id, 
			date_created, 
			oof_shard
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.DeliveryUID,
		order.PaymentTransaction,
		order.ItemsRID,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert orders data: %v\n", err)
		os.Exit(1)
	}
}

func DbInsertDeliveries(delivery Delivery) {
	conn := DbConnect(postgresUrl)
	defer conn.Close(context.Background())

	_, err := conn.Exec(context.Background(),
		`INSERT INTO deliveries (
			delivery_uid, name, phone, zip,
			city, address, region, email
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)`,
		delivery.DeliveryUID,
		delivery.Name,
		delivery.Phone,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert delivery data: %v\n", err)
		os.Exit(1)
	}
}

func DbInsertPayments(payment Payment) {
	conn := DbConnect(postgresUrl)
	defer conn.Close(context.Background())

	_, err := conn.Exec(context.Background(),
		`INSERT INTO payments (
			transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost,
			goods_total, custom_fee
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert payment data: %v\n", err)
		os.Exit(1)
	}
}

func DbInsertItems(items []Item) {
	conn := DbConnect(postgresUrl)
	defer conn.Close(context.Background())

	for _, item := range items {
		_, err := conn.Exec(context.Background(),
			`INSERT INTO items (
			chrt_id, track_number, price,
			rid, name, sale, size, total_price,
			nm_id, brand, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)`,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inserting item: %v", err)
		}
	}
}

func DbGetRowById(id string) (order Order, err error) {
	conn := DbConnect(postgresUrl)
	defer conn.Close(context.Background())

	//Get order row
	rows, _ := conn.Query(context.Background(), "select * from orders where order_uid = $1", id)
	orderData, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Order])
	if err != nil {
		panic("Collect order rows error")
	}

	delivery, _ := conn.Query(context.Background(), "select * from deliveries where delivery_uid = $1", id)
	orderData.Delivery, err = pgx.CollectOneRow(delivery, pgx.RowToStructByName[Delivery])

	payment, _ := conn.Query(context.Background(), "select * from payments where transaction = $1", id)
	orderData.Payment, err = pgx.CollectOneRow(payment, pgx.RowToStructByName[Payment])

	//Get every item RID
	for _, rid := range orderData.ItemsRID {
		fmt.Println(rid)
		items, _ := conn.Query(context.Background(), "select * from items where rid = $1", rid)
		dataItems, err := pgx.CollectOneRow(items, pgx.RowToStructByName[Item])
		if err != nil {
			panic("Failed to collect rows for items array")
		}

		orderData.Items = append(orderData.Items, dataItems)
	}

	return orderData, err
}
