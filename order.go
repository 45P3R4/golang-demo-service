package main

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderUID           uuid.UUID   `json:"order_uid" db:"order_uid"`
	TrackNumber        string      `json:"track_number" db:"track_number"`
	Entry              string      `json:"entry" db:"entry"`
	DeliveryUID        uuid.UUID   `db:"delivery_uid"`
	Delivery           Delivery    `json:"delivery,omitempty" db:"-"`
	PaymentTransaction uuid.UUID   `db:"payment_transaction"`
	Payment            Payment     `json:"payment,omitempty" db:"-"`
	ItemsRID           []uuid.UUID `db:"items_rid"`
	Items              []Item      `json:"items,omitempty" db:"-"`
	Locale             string      `json:"locale" db:"locale"`
	InternalSignature  string      `json:"internal_signature,omitempty" db:"internal_signature"`
	CustomerID         string      `json:"customer_id" db:"customer_id"`
	DeliveryService    string      `json:"delivery_service" db:"delivery_service"`
	Shardkey           int         `json:"shardkey" db:"shardkey"`
	SmID               int         `json:"sm_id" db:"sm_id"`
	DateCreated        time.Time   `json:"date_created" db:"date_created"`
	OofShard           int         `json:"oof_shard" db:"oof_shard"`
}

type Delivery struct {
	DeliveryUID uuid.UUID `json:"delivery_uid" db:"delivery_uid"`
	Name        string    `json:"name" db:"name"`
	Phone       string    `json:"phone" db:"phone"`
	Zip         int       `json:"zip" db:"zip"`
	City        string    `json:"city" db:"city"`
	Address     string    `json:"address" db:"address"`
	Region      string    `json:"region" db:"region"`
	Email       string    `json:"email" db:"email"`
}

type Payment struct {
	Transaction  uuid.UUID `json:"transaction" db:"transaction"`
	RequestID    string    `json:"request_id,omitempty" db:"request_id"`
	Currency     string    `json:"currency" db:"currency"`
	Provider     string    `json:"provider" db:"provider"`
	Amount       int       `json:"amount" db:"amount"`
	PaymentDt    int       `json:"payment_dt" db:"payment_dt"`
	Bank         string    `json:"bank" db:"bank"`
	DeliveryCost float64   `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int       `json:"goods_total" db:"goods_total"`
	CustomFee    int       `json:"custom_fee" db:"custom_fee"`
}

type Item struct {
	ChrtID      int       `json:"chrt_id" db:"chrt_id"`
	TrackNumber string    `json:"track_number" db:"track_number"`
	Price       int       `json:"price" db:"price"`
	RID         uuid.UUID `json:"rid" db:"rid"`
	Name        string    `json:"name" db:"name"`
	Sale        int       `json:"sale" db:"sale"`
	Size        int       `json:"size" db:"size"`
	TotalPrice  int       `json:"total_price" db:"total_price"`
	NmID        int       `json:"nm_id" db:"nm_id"`
	Brand       *string   `json:"brand,omitempty" db:"brand"`
	Status      int       `json:"status" db:"status"`
}
