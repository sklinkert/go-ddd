package response

import "time"

type ProductResponse struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	Currency   string    `json:"currency"`
	SellerId   string    `json:"seller_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ListProductsResponse struct {
	Products []*ProductResponse `json:"products"`
}
