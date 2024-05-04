package response

import "time"

type ProductResponse struct {
	Id        string
	Name      string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListProductsResponse struct {
	Products []*ProductResponse `json:"Products"`
}
