package request

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
	"mime/multipart"
)

type Product struct {
	Name             string                `json:"name" form:"name"`
	Description      string                `json:"description" form:"description"`
	Quantity         int                   `json:"quantity" form:"quantity"`
	UnitQuantity     string                `json:"unit_quantity" form:"unit_quantity"`
	Price            float64               `json:"price" form:"price"`
	UnitPrice        string                `json:"unit_price" form:"unit_price"`
	Image            string                `json:"-"`
	ImagePath        string                `json:"-"`
	Status           enum.StatusProduct    `json:"status" form:"status"`
	IsPreOrder       bool                  `json:"is_pre_order"  form:"is_pre_order"`
	MinPrice         float64               `json:"min_price" form:"min_price"`
	MaxPrice         float64               `json:"max_price" form:"max_price"`
	ProductCreatedAt string                `json:"product_created_at" form:"product_created_at"`
	ExpiredAt        string                `json:"expired_at" form:"expired_at"`
	CompanyID        int                   `json:"company_id" form:"company_id"`
	Commodity        string                `json:"commodity" form:"commodity"`
	File             *multipart.FileHeader `json:"-" form:"file"`
	IsActive         bool                  `json:"is_active" form:"is_active"`
	TmpImagePath     string                `json:"-"`
}

type ProductPaged struct {
	CompanyID int    `uri:"company_id"`
	Search    string `form:"search"`
	Page      int    `form:"page"`
	Size      int    `form:"size"`
	Commodity string `form:"commodity"`
	Status    string `form:"status"`
}

type ProductDetail struct {
	ID     int `uri:"id"`
	UserID int `form:"user_id"`
}
