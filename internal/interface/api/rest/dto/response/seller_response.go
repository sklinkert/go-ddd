package response

type SellerResponse struct {
	Id   string
	Name string
}

type ListSellersResponse struct {
	Sellers []*SellerResponse
}
