package model

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
	"time"
)

type Product struct {
	ID               int                `json:"id" gorm:"primary_key"`
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Quantity         int                `json:"quantity" form:"quantity"`
	UnitQuantity     string             `json:"unit_quantity"`
	Price            float64            `json:"price"`
	UnitPrice        string             `json:"unit_price"`
	Image            string             `json:"image"`
	TmpImagePath     string             `json:"-"`
	Status           enum.StatusProduct `json:"status"`
	ProductCreatedAt string             `json:"product_created_at"`
	ExpiredAt        string             `json:"expired_at"`
	Commodity        string             `json:"commodity"`
	CompanyID        int                `json:"company_id"`
	IsPreOrder       bool               `json:"is_pre_order"  gorm:"default:false"`
	MinPrice         float64            `json:"min_price"`
	MaxPrice         float64            `json:"max_price"`
	PariProductId    string             `json:"pari_product_id" form:"pari_product_id"`
	IsActive         bool               `json:"is_active" gorm:"default:true"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	DeletedAt        *time.Time         `sql:"index" json:"deleted_at"`
	Transaction      []PariTransaction  `json:"transaction" gorm:"-"`
}
