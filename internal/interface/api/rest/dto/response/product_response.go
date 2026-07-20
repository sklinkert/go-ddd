package response

import "time"

type ProductResponse struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	PriceMinorUnits int64     `json:"price_minor_units"`
	Currency        string    `json:"currency"`
	SellerId        string    `json:"seller_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ListProductsResponse struct {
	Products []*ProductResponse `json:"products"`
}
