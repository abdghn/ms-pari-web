/*
 * Created on 01/04/22 15.01
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package model

import (
	"bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"
	"time"
)

type User struct {
	ID                 int                    `json:"id" gorm:"primary_key"`
	Name               string                 `json:"name"`
	Email              string                 `json:"email"  gorm:"unique"`
	VerificationLevel  enum.VerificationLevel `sql:"verification_level"`
	Password           string                 `json:"password"`
	RoleID             int                    `json:"role_id" gorm:"column:role_id"`
	RoleName           string                 `json:"role_name" gorm:"-"`
	CompanyID          int                    `json:"company_id" gorm:"column:company_id"`
	CompanyName        string                 `json:"company_name" gorm:"-"`
	MustChangePassword bool                   `json:"must_change_password" gorm:"default:true"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	DeletedAt          *time.Time             `sql:"index" json:"deleted_at"`
}
