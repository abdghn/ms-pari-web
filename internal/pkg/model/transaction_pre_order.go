package model

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
	"time"
)

type TransactionPreOrder struct {
	ID                int                `json:"id" gorm:"primary_key"`
	PariProductID     string             `json:"pari_product_id"`
	PariTransactionID string             `json:"pari_transaction_id"`
	ProductID         int                `json:"product_id" gorm:"column:product_id"`
	ProductName       string             `json:"product_name,omitempty" gorm:"-"`
	ProductCommodity  string             `json:"product_commodity,omitempty" gorm:"-"`
	ProductImage      string             `json:"product_image,omitempty" gorm:"-"`
	ProductMinPrice   float64            `json:"product_min_price,omitempty" gorm:"-"`
	ProductMaxPrice   float64            `json:"product_max_price,omitempty" gorm:"-"`
	ProductExpiredAt  string             `json:"product_expired_at,omitempty" gorm:"-"`
	ProductCreatedAt  string             `json:"product_created_at,omitempty" gorm:"-"`
	ProductIsPreOrder bool               `json:"product_is_pre_order,omitempty" gorm:"-"`
	ProductIsActive   bool               `json:"product_is_active,omitempty" gorm:"-"`
	CompanyID         int                `json:"company_id" gorm:"column:company_id"`
	Quantity          int                `json:"quantity"`
	Status            enum.StatusProduct `json:"status"`
	ActualPrice       float64            `json:"actual_price"`
	BuyerName         string             `json:"buyer_name"`
	BuyerAddress      string             `json:"buyer_address"`
	BuyerContact      string             `json:"buyer_contact"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
	DeletedAt         *time.Time         `sql:"index" json:"deleted_at"`
}
