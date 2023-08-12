package postgres

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewConnection() (*gorm.DB, error) {
	return gorm.Open("postgres", "your_connection_string_here")
}
