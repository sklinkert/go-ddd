package rest

type CreateProductRequest struct {
	Name     string  `json:"Name"`
	Price    float64 `json:"Price"`
	SellerID string  `json:"SellerId"`
}
