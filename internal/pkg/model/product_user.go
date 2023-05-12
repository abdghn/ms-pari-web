/*
 * Created on 26/05/22 23.20
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package model

import (
	"time"
)

type ProductUser struct {
	ID        int        `json:"id" gorm:"primary_key"`
	ProductID int        `json:"product_id" gorm:"column:product_id"`
	UserID    int        `json:"user_id" gorm:"column:user_id"`
	CompanyID int        `json:"company_id" gorm:"column:company_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
