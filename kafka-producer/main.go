package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

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

func oneOf(options []string) string {
	return options[rand.Intn(len(options))]
}

func randomTime() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}

func generateRandomOrder() Order {
	rand.Seed(time.Now().UnixNano())

	order_id, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("Failed to create random order uuid", err)
	}

	// Генерируем Delivery
	delivery := Delivery{
		DeliveryUID: order_id,
		Name:        faker.Name(),
		Phone:       faker.Phonenumber(),
		Zip:         rand.Intn(99999) + 10000,
		City:        oneOf([]string{"Moscow", "Bobrov", "Vologda", "Vladivostok", "Saratov", "Omsk", "Kursk", "Voronezh", "Lipetsk"}),
		Address:     oneOf([]string{"Prospekt Mira 15", "Ploshad Mira 20", "Ulitsa Lenina 23", "Bobrovskaya 17", "Peredovoy 56", "Altayskaya 17"}),
		Region:      "Russia",
		Email:       faker.Email(),
	}

	// Генерируем Payment
	payment := Payment{
		Transaction:  order_id,
		RequestID:    "",
		Currency:     "RUB",
		Provider:     oneOf([]string{"wbpay", "sberpay", "tinkoff"}),
		Amount:       rand.Intn(10000) + 1000,
		PaymentDt:    int(time.Now().Unix()),
		Bank:         oneOf([]string{"sber", "alpha", "tinkoff"}),
		DeliveryCost: float64(rand.Intn(2000)+500) / 100,
		GoodsTotal:   rand.Intn(5000) + 100,
		CustomFee:    rand.Intn(100),
	}

	// Генерируем Items (3-5 товаров)
	itemCount := rand.Intn(3) + 3
	items := make([]Item, itemCount)
	itemsRID := make([]uuid.UUID, itemCount)

	brands := []string{"Nike", "Adidas", "Puma", "Reebok", "Apple", "Samsung", "Xiaomi", "Bosch"}

	for i := 0; i < itemCount; i++ {
		item_id, err := uuid.NewRandom()
		if err != nil {
			fmt.Printf("Failed to create random item uuid", err)
		}

		price := rand.Intn(5000) + 100
		sale := rand.Intn(70)
		brand := brands[rand.Intn(len(brands))]

		items[i] = Item{
			ChrtID:      rand.Intn(9999999) + 1000000,
			TrackNumber: fmt.Sprintf("WBILM%08d", rand.Intn(99999999)),
			Price:       price,
			RID:         item_id,
			Name:        faker.Word(),
			Sale:        sale,
			Size:        rand.Int() % 100,
			TotalPrice:  price * (100 - sale) / 100,
			NmID:        rand.Intn(9999999) + 1000000,
			Brand:       &brand,
			Status:      202,
		}
		itemsRID[i] = items[i].RID
	}

	// Собираем Order
	order := Order{
		OrderUID:           order_id,
		TrackNumber:        fmt.Sprintf("WBIL%012d", rand.Intn(999999999999)),
		Entry:              "WBIL",
		DeliveryUID:        delivery.DeliveryUID,
		Delivery:           delivery,
		PaymentTransaction: payment.Transaction,
		Payment:            payment,
		ItemsRID:           itemsRID,
		Items:              items,
		Locale:             oneOf([]string{"en", "ru", "jp", "it", "ge"}),
		InternalSignature:  "",
		CustomerID:         fmt.Sprintf("user%04d", rand.Intn(9999)),
		DeliveryService:    oneOf([]string{"meest", "pony", "dhl", "fedex"}),
		Shardkey:           rand.Intn(10) + 1,
		SmID:               rand.Intn(100) + 1,
		DateCreated:        randomTime(),
		OofShard:           rand.Intn(2) + 1,
	}

	return order
}

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

	order := generateRandomOrder()

	jsonData, err := json.Marshal(order)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: jsonData},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	fmt.Println("message created")

}

func produceInvalidData() {
	var conn, err = connect_kafka("orders", 0)
	if err != nil {
		fmt.Println("Failed to connect Kafka")
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	jsonData := []byte("test invalid data")

	_, err = conn.WriteMessages(
		kafka.Message{Value: jsonData},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	fmt.Println("message created")

}

func main() {

	produceTestData()
	// produceInvalidData()

}
