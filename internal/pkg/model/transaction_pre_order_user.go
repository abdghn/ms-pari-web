package model

import "time"

type TransactionPreOrderUser struct {
	ID                    int        `json:"id" gorm:"primary_key"`
	TransactionPreOrderID int        `json:"transaction_pre_order_id " gorm:"column:transaction_pre_order_id"`
	UserID                int        `json:"user_id" gorm:"column:user_id"`
	CompanyID             int        `json:"company_id" gorm:"column:company_id"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `sql:"index" json:"deleted_at"`
}
