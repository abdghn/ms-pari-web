package model

import (
	"time"
)

type PariProduct struct {
	ID               string    `json:"id"`
	ProductName      string    `json:"product_name"`
	ProductCommodity string    `json:"product_commodity"`
	Images           string    `json:"images"`
	Price            string    `json:"price"`
	CorporateID      string    `json:"corporate_id"`
	Status           string    `json:"status"`
	IsPreOrder       string    `json:"isPreOrder"`
	MinPrice         string    `json:"minPrice"`
	MaxPrice         string    `json:"maxPrice"`
	CreatedAt        time.Time `json:"created_at"`
}

type PariTransaction struct {
	IDProduct string    `json:"id_product"`
	IDBuyer   string    `json:"id_buyer"`
	Price     string    `json:"price"`
	Quantity  string    `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PariProductDetail struct {
	ID               string            `json:"id"`
	ProductName      string            `json:"product_name"`
	ProductCommodity string            `json:"product_commodity"`
	Images           string            `json:"images"`
	Price            string            `json:"price"`
	CorporateID      int               `json:"corporate_id"`
	Status           int               `json:"status"`
	IsPreOrder       int               `json:"isPreOrder"`
	MinPrice         int               `json:"minPrice"`
	MaxPrice         int               `json:"maxPrice"`
	CreatedAt        time.Time         `json:"created_at"`
	Transaction      []PariTransaction `json:"transaction"`
}
