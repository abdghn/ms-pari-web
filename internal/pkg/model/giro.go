package model

import "time"

type Giro struct {
	ID          int        `json:"id" gorm:"primary_key"`
	Code        string     `json:"code"  gorm:"unique"`
	CompanyName string     `json:"company_name"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at"`
}
