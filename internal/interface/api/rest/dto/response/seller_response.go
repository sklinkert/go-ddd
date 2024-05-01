package response

type SellerResponse struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
}

type ListSellersResponse struct {
	Sellers []*SellerResponse `json:"Sellers"`
}
