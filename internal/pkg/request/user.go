/*
 * Created on 07/04/22 05.24
 *
 * Copyright (c) 2022 Abdul Ghani Abbasi
 */

package request

import "bitbucket.org/bridce/ms-pari-web/internal/pkg/enum"

type User struct {
	Name               string                 `json:"name"`
	Email              string                 `json:"email"`
	Password           string                 `json:"password"`
	Role               string                 `json:"role"`
	CompanyID          int                    `json:"company_id"`
	VerificationLevel  enum.VerificationLevel `json:"verification_level"`
	MustChangePassword bool                   `json:"must_change_password"`
}

type Users []User

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Password string `json:"password"`
	UserID   int    `json:"-"`
}

type OpenKey struct {
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}
