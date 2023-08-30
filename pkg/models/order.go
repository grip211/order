package models

import (
	"time"

	"github.com/grip211/order/domain"
)

type Order struct {
	OrderUID          string          `db:"order_uid"`
	TrackNumber       string          `db:"track_number"`
	Entry             string          `db:"entry"`
	Delivery          domain.Delivery `db:"delivery"`
	Locale            string          `db:"locale"`
	InternalSignature string          `db:"internal_signature"`
	CustomerID        string          `db:"customer_id"`
	DeliveryService   string          `db:"delivery_service"`
	Shardkey          string          `db:"shardkey"`
	SmID              int             `db:"sm_id"`
	DateCreated       time.Time       `db:"date_created"`
	OofShard          string          `db:"oof_shard"`
}

func OrderModelsFromEntity(order *domain.Order) (Order, []Item, Payment) {
	return orderFromEntity(order), itemsFromEntity(order.Items), paymentFromEntity(order.Payment)
}

func orderFromEntity(order *domain.Order) Order {
	return Order{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          order.Delivery,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}
}

type Item struct {
	ChrtID      int    `db:"chrt_id"`
	TrackNumber string `db:"track_number"`
	Price       int    `db:"price"`
	Rid         string `db:"rid"`
	Name        string `db:"name"`
	Sale        int    `db:"sale"`
	Size        string `db:"size"`
	TotalPrice  int    `db:"total_price"`
	NmID        int    `db:"nm_id"`
	Brand       string `db:"brand"`
	Status      int    `db:"status"`
}

type Items []Item

func (i Items) toEntity() []domain.Item {
	models := make([]domain.Item, 0, len(i))
	for _, item := range i {
		models = append(models, domain.Item{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.ChrtID,
			Brand:       item.Brand,
			Status:      item.Status,
		})
	}
	return models
}

func itemsFromEntity(items []domain.Item) []Item {
	models := make([]Item, 0, len(items))
	for _, item := range items {
		models = append(models, Item{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.ChrtID,
			Brand:       item.Brand,
			Status:      item.Status,
		})
	}
	return models
}

type Payment struct {
	TransactionID string `db:"transaction_id"`
	RequestID     string `db:"request_id"`
	Currency      string `db:"currency"`
	Provider      string `db:"provider"`
	Amount        int    `db:"amount"`
	PaymentDt     int    `db:"payment_dt"`
	Bank          string `db:"bank"`
	DeliveryCost  int    `db:"delivery_cost"`
	GoodsTotal    int    `db:"goods_total"`
	CustomFee     int    `db:"custom_fee"`
}

func (p *Payment) toEntity() domain.Payment {
	return domain.Payment{
		Transaction:  p.TransactionID,
		RequestID:    p.RequestID,
		Currency:     p.Currency,
		Provider:     p.Provider,
		Amount:       p.Amount,
		PaymentDt:    p.PaymentDt,
		Bank:         p.Bank,
		DeliveryCost: p.DeliveryCost,
		GoodsTotal:   p.GoodsTotal,
		CustomFee:    p.CustomFee,
	}
}

func paymentFromEntity(payment domain.Payment) Payment {
	return Payment{
		TransactionID: payment.Transaction,
		RequestID:     payment.RequestID,
		Currency:      payment.Currency,
		Provider:      payment.Provider,
		Amount:        payment.Amount,
		PaymentDt:     payment.PaymentDt,
		Bank:          payment.Bank,
		DeliveryCost:  payment.DeliveryCost,
		GoodsTotal:    payment.GoodsTotal,
		CustomFee:     payment.CustomFee,
	}
}

type OrderReadModel struct {
	Order   Order   `db:"order"`
	Payment Payment `db:"payment"`
}

func (o *OrderReadModel) ToEntity(items Items) domain.Order {
	return domain.Order{
		OrderUID:          o.Order.OrderUID,
		TrackNumber:       o.Order.TrackNumber,
		Entry:             o.Order.Entry,
		Delivery:          o.Order.Delivery,
		Payment:           o.Payment.toEntity(),
		Items:             items.toEntity(),
		Locale:            o.Order.Locale,
		InternalSignature: o.Order.InternalSignature,
		CustomerID:        o.Order.CustomerID,
		DeliveryService:   o.Order.DeliveryService,
		Shardkey:          o.Order.Shardkey,
		SmID:              o.Order.SmID,
		DateCreated:       o.Order.DateCreated,
		OofShard:          o.Order.OofShard,
	}
}
