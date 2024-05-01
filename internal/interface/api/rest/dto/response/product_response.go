package response

type ProductResponse struct {
	Id    string
	Name  string
	Price float64
}

type ListProductsResponse struct {
	Products []*ProductResponse `json:"Products"`
}
