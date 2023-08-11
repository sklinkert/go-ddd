package entities

type Product struct {
	ID     uint
	Name   string
	Price  float64
	Seller *Seller
}
