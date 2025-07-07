package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

var dbConnection = DbConnect()

func DbConnect() *pgx.Conn {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	postgresURL := "postgresql://localhost/WBORDERS?user=" + dbUser + "&password=" + dbPassword

	conn, err := pgx.Connect(context.Background(), postgresURL)
	if err != nil {
		panic("[DbConnect]: Unable to connect to database: " + err.Error())
	}
	return conn
}

func DbInsert(m kafka.Message) {
	var dataOrder Order

	dataOrder.ItemsRID = make([]uuid.UUID, 0)
	dataOrder.Items = make([]Item, 0)

	err := json.Unmarshal(m.Value, &dataOrder)
	if err != nil {
		panic("[DbInsert]: failed to unmarshal JSON: " + err.Error())
	}

	DbInsertItems(dataOrder.Items)
	DbInsertDeliveries(dataOrder.Delivery)
	DbInsertPayments(dataOrder.Payment)

	dataOrder.DeliveryUID = dataOrder.Delivery.DeliveryUID
	dataOrder.PaymentTransaction = dataOrder.Payment.Transaction

	DbInsertOrders(dataOrder)
}

func DbInsertOrders(order Order) {
	_, err := dbConnection.Exec(context.Background(),
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
		panic("[DbInsertOrders]: Unable to insert orders data: " + err.Error())
	}
}

func DbInsertDeliveries(delivery Delivery) {

	_, err := dbConnection.Exec(context.Background(),
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
		panic("[DbInsertDeliveries]: Unable to insert delivery data: " + err.Error())
	}
}

func DbInsertPayments(payment Payment) {

	_, err := dbConnection.Exec(context.Background(),
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
		panic("[DbInsertPayments]: Unable to insert payment data: " + err.Error())
	}
}

func DbInsertItems(items []Item) {

	for _, item := range items {
		_, err := dbConnection.Exec(context.Background(),
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
			panic("[DbInsertItems]: Unable to insert item data: " + err.Error())
		}
	}
}

func DbGetRowById(id string) (order Order, err error) {
	//Get order row
	rows, _ := dbConnection.Query(context.Background(), "select * from orders where order_uid = $1", id)
	orderData, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Order])
	if err != nil {
		panic("[DbGetRowById]: Failed to collect order row:" + err.Error())
	}

	//Get delivery by order id
	delivery, _ := dbConnection.Query(context.Background(), "select * from deliveries where delivery_uid = $1", id)
	orderData.Delivery, err = pgx.CollectOneRow(delivery, pgx.RowToStructByName[Delivery])
	if err != nil {
		panic("[DbGetRowById]: Failed to collect delivery row: " + err.Error())
	}

	//Get payment by order id
	payment, _ := dbConnection.Query(context.Background(), "select * from payments where transaction = $1", id)
	orderData.Payment, err = pgx.CollectOneRow(payment, pgx.RowToStructByName[Payment])
	if err != nil {
		panic("[DbGetRowById]: Failed to collect payment row:" + err.Error())
	}

	//Get every item RID
	for _, rid := range orderData.ItemsRID {
		items, _ := dbConnection.Query(context.Background(), "select * from items where rid = $1", rid)
		dataItems, err := pgx.CollectOneRow(items, pgx.RowToStructByName[Item])
		if err != nil {
			panic("[DbGetRowById]: Failed to collect rows for items array:" + err.Error())
		}

		orderData.Items = append(orderData.Items, dataItems)
	}

	return orderData, err
}

func DbGetLastRows(count int) (order []Order, err error) {

	//Get order row
	rows, _ := dbConnection.Query(context.Background(), "select * from orders order by date_created limit $1", count)
	orderData, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Order])
	if err != nil {
		panic("[DbGetLastRows]: Failder to collect rows: " + err.Error())
	}

	for i := range orderData {
		order := &orderData[i]

		//Get delivery by order id
		delivery, _ := dbConnection.Query(context.Background(), "select * from deliveries where delivery_uid = $1 ", order.OrderUID)
		order.Delivery, err = pgx.CollectOneRow(delivery, pgx.RowToStructByName[Delivery])
		if err != nil {
			panic("[DbGetLastRows]: Failed to collect delivery row:" + err.Error())
		}

		//Get payment by order id
		payment, _ := dbConnection.Query(context.Background(), "select * from payments where transaction = $1", order.OrderUID)
		order.Payment, err = pgx.CollectOneRow(payment, pgx.RowToStructByName[Payment])
		if err != nil {
			panic("[DbGetLastRows]: Failed to collect payment row:" + err.Error())
		}

		//Get every item RID
		for _, rid := range order.ItemsRID {
			items, _ := dbConnection.Query(context.Background(), "select * from items where rid = $1", rid)
			dataItems, err := pgx.CollectOneRow(items, pgx.RowToStructByName[Item])
			if err != nil {
				panic("[DbGetLastRows]: Failed to collect item row: " + err.Error())
			}
			order.Items = append(order.Items, dataItems)
		}
	}

	return orderData, err
}
