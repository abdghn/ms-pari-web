package model

import "time"

type Company struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"name"  gorm:"unique"`
	Code      string     `json:"code"  gorm:"unique"`
	Alias     string     `json:"alias"`
	Address   string     `json:"address"`
	Giro      string     `json:"giro"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
