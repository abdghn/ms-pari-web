package request

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
)

type TransactionPreOrder struct {
	PariProductId     string             `json:"pari_product_id" form:"pari_product_id"`
	PariTransactionId string             `json:"pari_transaction_id" form:"pari_transaction_id"`
	ProductID         int                `json:"product_id" form:"product_id"`
	CompanyID         int                `json:"company_id" form:"company_id"`
	Quantity          int                `json:"quantity" form:"quantity"`
	BuyerName         string             `json:"buyer_name" form:"buyer_name"`
	BuyerAddress      string             `json:"buyer_address" form:"buyer_address"`
	BuyerContact      string             `json:"buyer_contact" form:"buyer_contact"`
	ActualPrice       float64            `json:"actual_price" form:"actual_price"`
	Status            enum.StatusProduct `json:"status" form:"status"`
}

type TransactionPreOrderPaged struct {
	CompanyID int    `uri:"company_id"`
	Search    string `form:"search"`
	Page      int    `form:"page"`
	Size      int    `form:"size"`
	Commodity string `form:"commodity"`
	Status    string `form:"status"`
}

type TransactionPreOrderDetail struct {
	ID     int `uri:"id"`
	UserID int `form:"user_id"`
}
