/*
 * Created on 07/04/22 08.43
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package model

import "time"

type Role struct {
	ID        int        `json:"id" gorm:"primary_key"`
	Name      string     `json:"name"  gorm:"unique"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}
