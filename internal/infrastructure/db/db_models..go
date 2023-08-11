package db

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Name     string
	Price    float64
	SellerID uint
}

type Seller struct {
	gorm.Model
	Name string
}
