package response

import "time"

type SellerResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListSellersResponse struct {
	Sellers []*SellerResponse `json:"sellers"`
}
