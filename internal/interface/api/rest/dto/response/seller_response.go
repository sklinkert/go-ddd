package response

import "time"

type SellerResponse struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListSellersResponse struct {
	Sellers []*SellerResponse
}
